package dynamodb

import (
	"email_action/logging"
	"email_action/models"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

var (
	log                   = logging.NewZapLogger()
	itemNotFoundErr       = "Item not found"
	ErrItemNotFound       = errors.New("item not found")
	ErrMismatchedPassword = errors.New("incorrect password")
)

type DynamodbAdapter struct {
	dynamodbSvc        dynamodbiface.DynamoDBAPI
	usersTable         *string
	projectsTable      *string
	projectsTableIndex *string
	articlesTable      *string
	subscribersTable   *string
	resetPasswordTable *string
	ticketTable        *string
}

func NewDynamodbAdapter(env string) *DynamodbAdapter {
	log.Infof("NewDynamodbAdapter(): env = %s", env)
	awsConfig := &aws.Config{
		Region: aws.String("us-west-2"),
	}
	sess := session.Must(session.NewSession(awsConfig))
	return &DynamodbAdapter{
		dynamodbSvc:        dynamodb.New(sess, awsConfig),
		usersTable:         aws.String("users-" + env),
		projectsTable:      aws.String("projects-" + env),
		articlesTable:      aws.String("articles-" + env),
		subscribersTable:   aws.String("subscribers-" + env),
		projectsTableIndex: aws.String("projects-index-projectId-" + env),
		resetPasswordTable: aws.String("reset-password-" + env),
		ticketTable:        aws.String("tickets-" + env),
	}
}

func Marshal(in interface{}) (*dynamodb.AttributeValue, error) {
	encoder := dynamodbattribute.NewEncoder(func(encoder *dynamodbattribute.Encoder) {
		encoder.EnableEmptyCollections = true
	})
	return encoder.Encode(in)
}

func MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	encoder := dynamodbattribute.NewEncoder(func(encoder *dynamodbattribute.Encoder) {
		encoder.EnableEmptyCollections = true
	})
	av, err := encoder.Encode(in)
	if err != nil || av == nil || av.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}
	return av.M, nil
}

func (d *DynamodbAdapter) CreateUserIfNotExist(user *models.User) error {
	log.Infof("CreateUserIfNotExist():")
	tableItem, err := MarshalMap(user)
	if err != nil {
		log.Errorf("CreateUserIfNotExist(): failed to marshal user")
	}
	putItemInput := dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(email)"),
		Item:                tableItem,
		TableName:           d.usersTable,
	}
	_, err = d.dynamodbSvc.PutItem(&putItemInput)
	if err != nil {
		log.Errorf("CreateUserIfNotExist(): failed to put item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) ReadUser(email string) (*models.User, error) {
	log.Infof("ReadUser(): %s", email)
	getItemInput := dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
		},
		TableName: d.usersTable,
	}
	dbOutput, err := d.dynamodbSvc.GetItem(&getItemInput)
	if err != nil {
		log.Errorf("ReadUser(): failed to get item, error: %+v", err)
		return nil, err
	}
	if dbOutput == nil || dbOutput.Item == nil {
		err = fmt.Errorf(itemNotFoundErr)
		return nil, err
	}
	user := &models.User{}
	err = dynamodbattribute.UnmarshalMap(dbOutput.Item, user)
	if err != nil {
		log.Errorf("ReadUser(): failed to unmarshall, error: %+v", err)
		return nil, err
	}
	return user, nil
}

func (d *DynamodbAdapter) UpdateUserBasicInfo(userInfo *models.UserInfo) error {
	log.Infof("UpdateUserBasicInfo():")
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(userInfo.Email),
			},
		},
		TableName: d.usersTable,
		UpdateExpression: aws.String("SET firstName = :firstName, lastName = :lastName, address = :address," +
			" phoneNumber = :phoneNumber"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":firstName": {
				S: aws.String(userInfo.FirstName),
			},
			":lastName": {
				S: aws.String(userInfo.LastName),
			},
			":address": {
				S: aws.String(userInfo.Address),
			},
			":phoneNumber": {
				S: aws.String(userInfo.PhoneNumber),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdateUserBasicInfo(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) AuthUser(authReq *models.AuthRequest) (*models.Claims, error) {
	log.Infof("AuthUser():")
	user, err := d.ReadUser(authReq.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(authReq.Password)); err != nil {
		return nil, err
	}

	claims := &models.Claims{
		Email:     user.Email,
		UserScope: user.UserScope,
		Exp:       time.Now().AddDate(0, 0, 14).Unix(),
		Iat:       time.Now().Unix(),
	}
	return claims, nil
}

func (d *DynamodbAdapter) UpdateUserSubscription(email, subscriptionEndTimeMs string, isSubscribed bool) error {
	log.Infof("UpdateUserSubscription(%s, %s, %v):", email, subscriptionEndTimeMs, isSubscribed)
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
		},
		UpdateExpression: aws.String("SET isSubscribed = :isSubscribed SET subscriptionEndTimeMs = :subscriptionEndTimeMs"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isSubscribed": {
				BOOL: aws.Bool(isSubscribed),
			},
			":subscriptionEndTimeMs": {
				N: aws.String(subscriptionEndTimeMs),
			},
		},
		TableName: d.usersTable,
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdateUserSubscription(): failed to put item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) CreateProjectIfNotExist(project *models.Project) error {
	log.Infof("CreateProjectIfNotExist(): project = %v", project)
	tableItem, err := MarshalMap(project)
	if err != nil {
		log.Errorf("CreateProjectIfNotExist(): failed to marshal user")
	}
	userIdExistsCondition := dynamodb.Update{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(project.Email)},
		},
		ConditionExpression: aws.String("attribute_exists(email) and isBlock=:bool"),
		UpdateExpression:    aws.String("SET projectIds = list_append(projectIds, :projectId)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":projectId": {
				L: []*dynamodb.AttributeValue{
					{
						S: aws.String(project.ProjectId),
					},
				},
			},
			":bool": {
				BOOL: aws.Bool(false),
			},
		},
		TableName: d.usersTable,
	}
	putItemInput := dynamodb.Put{
		ConditionExpression: aws.String("attribute_not_exists(email) and attribute_not_exists(projectId)"),
		Item:                tableItem,
		TableName:           d.projectsTable,
	}

	transactionItem1 := dynamodb.TransactWriteItem{
		Update: &userIdExistsCondition,
	}
	transactionItem2 := dynamodb.TransactWriteItem{
		Put: &putItemInput,
	}
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{&transactionItem1, &transactionItem2},
	}

	_, err = d.dynamodbSvc.TransactWriteItems(&transactionInput)
	if err != nil {
		log.Errorf("CreateProjectIfNotExist(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) ReadProject(email, projectId string) (*models.Project, error) {
	log.Infof("ReadProject(): email = %s, projectId = %s", email, projectId)
	getItemInput := dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(email)},
			"projectId": {S: aws.String(projectId)},
		},
		TableName: d.projectsTable,
	}
	dbOutput, err := d.dynamodbSvc.GetItem(&getItemInput)
	if err != nil {
		log.Errorf("ReadProject(): failed to get item, error: %+v", err)
		return nil, err
	}
	if dbOutput == nil || dbOutput.Item == nil {
		err = fmt.Errorf(itemNotFoundErr)
		return nil, err
	}
	project := &models.Project{}
	err = dynamodbattribute.UnmarshalMap(dbOutput.Item, project)
	if err != nil {
		log.Errorf("ReadProject(): failed to unmarshall, error: %+v", err)
		return nil, err
	}
	return project, nil
}

func (d *DynamodbAdapter) ReadProjectBrief(projectId string) (*models.ProjectBrief, error) {
	log.Infof("ReadProjectBrief(): projectId = %s", projectId)
	expr, err := compileKeyConditionExpression(projectId)
	if err != nil {
		return nil, err
	}
	queryItemInput := dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 d.projectsTable,
		IndexName:                 d.projectsTableIndex,
	}
	dbOutput, err := d.dynamodbSvc.Query(&queryItemInput)
	if err != nil {
		log.Errorf("ReadProjectBrief(): failed to get item, error: %+v", err)
		return nil, err
	}
	if dbOutput == nil || len(dbOutput.Items) == 0 {
		err = fmt.Errorf(itemNotFoundErr)
		return nil, err
	}
	project := &models.Project{}
	err = dynamodbattribute.UnmarshalMap(dbOutput.Items[0], project)
	if err != nil {
		log.Errorf("ReadProjectBrief(): failed to unmarshall, error: %+v", err)
		return nil, err
	}
	projectBrief := &models.ProjectBrief{
		ProjectId:           project.ProjectId,
		Name:                project.Name,
		CreatedAt:           project.CreatedAt,
		Interval:            project.Interval,
		LastBroadcastTimeMs: project.LastBroadcastTimeMs,
		Intro:               project.Intro,
		BackgroundImageUrl:  project.BackgroundImageUrl,
		AvatarUrl:           project.AvatarUrl,
		OutgoingEmail:       project.OutgoingEmail,
		Author:              project.Author,
		TotalArticles:       len(project.ArticleIds),
		SubscriptionType:    project.SubscriptionType,
	}
	return projectBrief, nil
}

func compileKeyConditionExpression(projectId string) (expression.Expression, error) {
	var builder expression.Builder
	keyCondition := expression.Key("projectId").Equal(expression.Value(projectId))
	builder = builder.WithKeyCondition(keyCondition)
	return builder.Build()
}

// dynamodb limitation: idPairs length must be <= 100 and total response size <= 1MB
func (d *DynamodbAdapter) BatchReadProjects(idPairs [][]string) ([]*models.Project, error) {
	log.Infof("BatchReadProjects(): %v", idPairs)
	var keys []map[string]*dynamodb.AttributeValue
	for _, itemPair := range idPairs {
		email := itemPair[0]
		projectId := itemPair[1]
		key := map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"projectId": {
				S: aws.String(projectId),
			},
		}
		keys = append(keys, key)
	}
	batchGetItemInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			*d.projectsTable: {
				Keys: keys,
			},
		},
	}
	dbOutput, err := d.dynamodbSvc.BatchGetItem(batchGetItemInput)
	if err != nil {
		log.Errorf("BatchReadProjects(): failed to get item, error: %+v", err)
		return nil, err
	}

	var projects []*models.Project
	for _, dbItem := range dbOutput.Responses {
		for _, item := range dbItem {
			project := &models.Project{}
			err = dynamodbattribute.UnmarshalMap(item, project)
			if err != nil {
				err := fmt.Errorf("BatchReadProjects(): error = %v", err)
				return nil, err
			}
			projects = append(projects, project)
		}
		break
	}

	return projects, nil
}

func (d *DynamodbAdapter) SaveProject(project *models.Project) error {
	log.Infof("SaveProject(): project = %v", project)
	existingProject, err := d.ReadProject(project.Email, project.ProjectId)
	if err != nil {
		return err
	}
	if isListChanged(existingProject.ArticleIds, project.ArticleIds) {
		err = fmt.Errorf("updating articleIds directly is not allowed, use article create/delete api instead")
		return err
	}
	tableItem, err := MarshalMap(project)
	if err != nil {
		log.Errorf("SaveProject(): failed to marshal user")
	}
	if len(project.ArticleIds) == 0 {
		delete(tableItem, "articleIds")
	}

	userIdExistsCondition := dynamodb.ConditionCheck{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(project.Email)},
		},
		ConditionExpression: aws.String("attribute_exists(email) and isBlock=:bool"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bool": {
				BOOL: aws.Bool(false),
			},
		},
		TableName: d.usersTable,
	}

	putItemInput := dynamodb.Put{
		Item:                tableItem,
		ConditionExpression: aws.String("attribute_exists(email) and attribute_exists(projectId)"),
		TableName:           d.projectsTable,
	}

	transactionItem1 := dynamodb.TransactWriteItem{
		ConditionCheck: &userIdExistsCondition,
	}
	transactionItem2 := dynamodb.TransactWriteItem{
		Put: &putItemInput,
	}
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{&transactionItem1, &transactionItem2},
	}

	_, err = d.dynamodbSvc.TransactWriteItems(&transactionInput)
	if err != nil {
		log.Errorf("CreateProjectIfNotExist(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) DeleteProject(email, projectId string) error {
	log.Infof("DeleteProject(): email = %s, projectId = %s", email, projectId)
	project, err := d.ReadProject(email, projectId)
	if err != nil {
		if err.Error() == itemNotFoundErr {
			log.Warnf("DeleteProject(): key %s, %s not found", email, projectId)
			return nil
		}
		return err
	}
	for _, articleId := range project.ArticleIds {
		// not using transaction because of the 25 items limit imposed by dynamodb
		// assume each project contains more than 25 articles
		d.DeleteArticle(email, projectId, articleId)
	}

	var lastEvaluatedKey map[string]*dynamodb.AttributeValue
	for {
		queryInput := dynamodb.QueryInput{
			KeyConditions: map[string]*dynamodb.Condition{
				"projectId": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dynamodb.AttributeValue{
						{
							S: aws.String(projectId),
						},
					},
				},
			},
			TableName: d.subscribersTable,
		}
		if len(lastEvaluatedKey) > 0 {
			queryInput.ExclusiveStartKey = lastEvaluatedKey
		}
		dbOutput, err := d.dynamodbSvc.Query(&queryInput)
		if err != nil {
			log.Errorf("DeleteProject() error = %v", err)
		}

		subscriber := &models.Subscriber{}
		for _, item := range dbOutput.Items {
			err = dynamodbattribute.UnmarshalMap(item, subscriber)
			if err != nil {
				log.Errorf("DeleteProject(): failed to unmarshall, error: %+v", err)
				continue
			}
			// not using transaction because of the 25 items limit imposed by dynamodb
			// assume each project contains more than 25 subscribers
			d.DeleteSubscriber(projectId, subscriber.Email)
		}

		if dbOutput.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = dbOutput.LastEvaluatedKey
	}

	user, err := d.ReadUser(email)
	if err != nil {
		log.Errorf("DeleteProject(): error = %+v", err)
		return err
	}
	deletingIndex := 0
	for deletingIndex = range user.ProjectIds {
		if user.ProjectIds[deletingIndex] == projectId {
			break
		}
	}

	deleteItem := dynamodb.Delete{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(email)},
			"projectId": {S: aws.String(projectId)},
		},
		TableName: d.projectsTable,
	}
	updateProject := dynamodb.Update{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
		},
		ConditionExpression: aws.String("attribute_exists(email)"),
		UpdateExpression:    aws.String("REMOVE projectIds[" + strconv.Itoa(deletingIndex) + "]"),
		TableName:           d.usersTable,
	}

	transactionItem1 := dynamodb.TransactWriteItem{
		Delete: &deleteItem,
	}
	transactionItem2 := dynamodb.TransactWriteItem{
		Update: &updateProject,
	}
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{&transactionItem1, &transactionItem2},
	}

	_, err = d.dynamodbSvc.TransactWriteItems(&transactionInput)
	if err != nil {
		log.Errorf("DeleteProject(): error: %+v", err)
		return err
	}
	return nil

	return nil
}

func (d *DynamodbAdapter) ScanProject(limit int, exclusiveStartKey map[string]*dynamodb.AttributeValue) (
	[]*models.Project, map[string]*dynamodb.AttributeValue, error) {
	log.Infof("ScanProject(): limit = %d, exclusiveStartKey = %v", limit, exclusiveStartKey)
	scanInput := dynamodb.ScanInput{
		FilterExpression: aws.String("projectState = :projectState"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":projectState": {
				N: aws.String(strconv.Itoa(models.ProjectStateLive)),
			},
		},
		Limit:     aws.Int64(int64(limit)),
		TableName: d.projectsTable,
	}
	if len(exclusiveStartKey) > 0 {
		scanInput.ExclusiveStartKey = exclusiveStartKey
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanInput)
	if err != nil {
		log.Errorf("ScanProject(): error = %v", err)
	}

	var res []*models.Project
	for _, item := range dbOutput.Items {
		project := &models.Project{}
		err = dynamodbattribute.UnmarshalMap(item, project)
		if err != nil {
			log.Warnf("ScanProject(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, project)
	}
	return res, dbOutput.LastEvaluatedKey, nil
}

func (d *DynamodbAdapter) ScanUser(limit int, exclusiveStartKey map[string]*dynamodb.AttributeValue) (
	[]*models.UserAccountInfo, map[string]*dynamodb.AttributeValue, error) {
	log.Infof("ScanUser(): limit = %d, exclusiveStartKey = %v", limit, exclusiveStartKey)
	scanInput := dynamodb.ScanInput{
		Limit:     aws.Int64(int64(limit)),
		TableName: d.usersTable,
	}
	if len(exclusiveStartKey) > 0 {
		scanInput.ExclusiveStartKey = exclusiveStartKey
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanInput)
	if err != nil {
		log.Errorf("ScanUser(): error = %v", err)
	}

	var res []*models.UserAccountInfo
	for _, item := range dbOutput.Items {
		userAccountInfo := &models.UserAccountInfo{}
		err = dynamodbattribute.UnmarshalMap(item, userAccountInfo)
		if err != nil {
			log.Warnf("ScanUser(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, userAccountInfo)
	}
	return res, dbOutput.LastEvaluatedKey, nil
}

func (d *DynamodbAdapter) ClearEmailQuota(email string) error {
	log.Infof("ClearEmailQuota(): %s", email)
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName:        d.usersTable,
		UpdateExpression: aws.String("SET lastClearCycleTime = :lastClearCycleTime, emailUsageInCycle = :emailUsageInCycle"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":lastClearCycleTime": {
				N: aws.String(strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)),
			},
			":emailUsageInCycle": {
				N: aws.String(strconv.FormatInt(int64(0), 10)),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("ClearEmailQuota(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) UpdateProjectBroadcastStat(project *models.Project) error {
	log.Infof("UpdateProjectBroadcastStat(): project = %v", project)
	updateInput := dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(project.Email)},
			"projectId": {S: aws.String(project.ProjectId)},
		},
		UpdateExpression: aws.String("SET totalBroadcastCount = :totalBroadcastCount, " +
			"lastBroadcastCount = :lastBroadcastCount, lastBroadcastDuration = :lastBroadcastDuration," +
			"lastBroadcastTimeMs = :lastBroadcastTimeMs"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":totalBroadcastCount": {
				N: aws.String(strconv.FormatInt(project.TotalBroadcastCount, 10)),
			},
			":lastBroadcastCount": {
				N: aws.String(strconv.FormatInt(project.LastBroadcastCount, 10)),
			},
			":lastBroadcastDuration": {
				N: aws.String(strconv.FormatInt(project.LastBroadcastDuration, 10)),
			},
			":lastBroadcastTimeMs": {
				N: aws.String(strconv.FormatInt(project.LastBroadcastTimeMs, 10)),
			},
		},
		TableName: d.projectsTable,
	}

	_, err := d.dynamodbSvc.UpdateItem(&updateInput)
	if err != nil {
		log.Errorf("UpdateProjectBroadcastStat(): error = %v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) SaveArticle(article *models.Article) error {
	log.Infof("SaveArticle(): article = %v", article)
	tableItem, err := MarshalMap(article)
	if err != nil {
		log.Errorf("SaveArticle(): failed to marshal article")
	}

	userIdExistsCondition := dynamodb.ConditionCheck{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(article.Email)},
		},
		ConditionExpression: aws.String("attribute_exists(email) and isBlock=:bool"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bool": {
				BOOL: aws.Bool(false),
			},
		},
		TableName: d.usersTable,
	}
	putArticle := dynamodb.Put{
		Item:      tableItem,
		TableName: d.articlesTable,
	}
	updateProject := dynamodb.Update{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(article.Email)},
			"projectId": {S: aws.String(article.ProjectId)},
		},
		ConditionExpression: aws.String("attribute_exists(email) and attribute_exists(projectId)"),
		UpdateExpression:    aws.String("SET articleIds = list_append(if_not_exists(articleIds, :empty_list), :articleId)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":articleId": {
				L: []*dynamodb.AttributeValue{
					{
						S: aws.String(article.ArticleId),
					},
				},
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: d.projectsTable,
	}

	transactionItem1 := dynamodb.TransactWriteItem{
		ConditionCheck: &userIdExistsCondition,
	}
	transactionItem2 := dynamodb.TransactWriteItem{
		Put: &putArticle,
	}
	transactionItem3 := dynamodb.TransactWriteItem{
		Update: &updateProject,
	}
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{&transactionItem1, &transactionItem2, &transactionItem3},
	}

	_, err = d.dynamodbSvc.TransactWriteItems(&transactionInput)
	if err != nil {
		log.Errorf("SaveArticle(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) UpdateArticle(article *models.Article) error {
	log.Infof("UpdateArticle(): article = %v", article)
	tableItem, err := MarshalMap(article)
	if err != nil {
		log.Errorf("UpdateArticle(): failed to marshal article")
	}

	userIdExistsCondition := dynamodb.ConditionCheck{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(article.Email)},
		},
		ConditionExpression: aws.String("attribute_exists(email) and isBlock=:bool"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bool": {
				BOOL: aws.Bool(false),
			},
		},
		TableName: d.usersTable,
	}

	projectIdExistsCondition := dynamodb.ConditionCheck{
		Key: map[string]*dynamodb.AttributeValue{
			"projectId": {S: aws.String(article.ProjectId)},
			"email":     {S: aws.String(article.Email)},
		},
		ConditionExpression: aws.String("attribute_exists(projectId)"),
		TableName:           d.projectsTable,
	}
	putArticle := dynamodb.Put{
		Item:      tableItem,
		TableName: d.articlesTable,
	}

	transactionItem1 := dynamodb.TransactWriteItem{
		ConditionCheck: &userIdExistsCondition,
	}
	transactionItem2 := dynamodb.TransactWriteItem{
		ConditionCheck: &projectIdExistsCondition,
	}
	transactionItem3 := dynamodb.TransactWriteItem{
		Put: &putArticle,
	}

	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{&transactionItem1, &transactionItem2, &transactionItem3},
	}

	_, err = d.dynamodbSvc.TransactWriteItems(&transactionInput)
	if err != nil {
		log.Errorf("UpdateArticle(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) ReadArticle(email, articleId string) (*models.Article, error) {
	log.Infof("ReadArticle(): email = %s, projectId = %s", email, articleId)
	getItemInput := dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(email)},
			"articleId": {S: aws.String(articleId)},
		},
		TableName: d.articlesTable,
	}
	dbOutput, err := d.dynamodbSvc.GetItem(&getItemInput)
	if err != nil {
		log.Errorf("ReadArticle(): failed to get item, error: %+v", err)
		return nil, err
	}
	if dbOutput == nil || dbOutput.Item == nil {
		err = fmt.Errorf("item(%s, %s) not found", email, articleId)
		return nil, err
	}
	article := &models.Article{}
	err = dynamodbattribute.UnmarshalMap(dbOutput.Item, article)
	if err != nil {
		log.Errorf("ReadArticle(): failed to unmarshall, error: %+v", err)
		return nil, err
	}
	return article, nil
}

// dynamodb limitation: idPairs length must be <= 100 and total response size <= 1MB
func (d *DynamodbAdapter) BatchReadArticles(idPairs [][]string) ([]*models.Article, error) {
	log.Infof("BatchReadArticles(): %v", idPairs)
	var keys []map[string]*dynamodb.AttributeValue
	for _, itemPair := range idPairs {
		email := itemPair[0]
		projectId := itemPair[1]
		key := map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"articleId": {
				S: aws.String(projectId),
			},
		}
		keys = append(keys, key)
	}
	batchGetItemInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			*d.articlesTable: {
				Keys: keys,
			},
		},
	}
	dbOutput, err := d.dynamodbSvc.BatchGetItem(batchGetItemInput)
	if err != nil {
		log.Errorf("BatchReadArticles(): failed to get item, error: %+v", err)
		return nil, err
	}

	var articles []*models.Article
	for _, dbItem := range dbOutput.Responses {
		for _, item := range dbItem {
			article := &models.Article{}
			err = dynamodbattribute.UnmarshalMap(item, article)
			if err != nil {
				err := fmt.Errorf("BatchReadArticles(): error = %v", err)
				return nil, err
			}
			articles = append(articles, article)
		}
		break
	}

	return articles, nil
}

func (d *DynamodbAdapter) DeleteArticle(email, projectId, articleId string) error {
	log.Infof("DeleteArticle(): email = %s, projectId = %s, articleId = %s", email, projectId, articleId)
	project, err := d.ReadProject(email, projectId)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Warnf("DeleteArticle(): key %s, %s not found", email, articleId)
				return nil
			}
		}
		return err
	}
	deletingIndex := 0
	for deletingIndex = range project.ArticleIds {
		if project.ArticleIds[deletingIndex] == articleId {
			break
		}
	}
	deleteItem := dynamodb.Delete{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(email)},
			"articleId": {S: aws.String(articleId)},
		},
		TableName: d.articlesTable,
	}
	updateProject := dynamodb.Update{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(email)},
			"projectId": {S: aws.String(projectId)},
		},
		ConditionExpression: aws.String("attribute_exists(email) and attribute_exists(projectId)"),
		UpdateExpression:    aws.String("REMOVE articleIds[" + strconv.Itoa(deletingIndex) + "]"),
		TableName:           d.projectsTable,
	}

	transactionItem1 := dynamodb.TransactWriteItem{
		Delete: &deleteItem,
	}
	transactionItem2 := dynamodb.TransactWriteItem{
		Update: &updateProject,
	}
	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{&transactionItem1, &transactionItem2},
	}

	_, err = d.dynamodbSvc.TransactWriteItems(&transactionInput)
	if err != nil {
		log.Errorf("DeleteArticle(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) CreateSubscriberIfNotExist(subscriber *models.Subscriber) error {
	log.Infof("CreateSubscriberIfNotExist(): subscriber = %v", subscriber)
	projectBrief, err := d.ReadProjectBrief(subscriber.ProjectId)
	if err != nil {
		return err
	}
	if projectBrief.SubscriptionType == models.RollingSubscription {
		subscriber.ArticleCursor = projectBrief.TotalArticles
	}

	tableItem, err := MarshalMap(subscriber)
	if err != nil {
		log.Errorf("CreateSubscriberIfNotExist(): failed to marshal user")
		return err
	}

	putItemInput := dynamodb.PutItemInput{
		Item:                tableItem,
		TableName:           d.subscribersTable,
		ConditionExpression: aws.String("attribute_not_exists(email) and attribute_not_exists(projectId)"),
	}
	_, err = d.dynamodbSvc.PutItem(&putItemInput)
	if err != nil {
		log.Debugf("CreateSubscriberIfNotExist(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) UpdateSubscriberBasic(req *models.UpdateSubscriberRequest) error {
	log.Infof("UpdateSubscriberBasic(): req = %v", req)
	updateItemInput := dynamodb.UpdateItemInput{
		TableName: d.subscribersTable,
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(req.Email),
			},
			"projectId": {
				S: aws.String(req.ProjectId),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":firstName": {
				S: aws.String(req.FirstName),
			},
			":lastName": {
				S: aws.String(req.LastName),
			},
			":isEnabled": {
				BOOL: aws.Bool(req.IsEnabled),
			},
			":articleCursor": {
				N: aws.String(strconv.Itoa(req.ArticleCursor)),
			},
		},
		UpdateExpression: aws.String("SET firstName = :firstName, lastName = :lastName," +
			" isEnabled = :isEnabled, articleCursor = :articleCursor"),
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("SaveSubscriber(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) UpdateSubscriberStat(subscriber *models.Subscriber) error {
	log.Infof("UpdateSubscriberStat(): subscriber = %v", subscriber)
	updateItemInput := dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(subscriber.Email)},
			"projectId": {S: aws.String(subscriber.ProjectId)},
		},
		TableName:        d.subscribersTable,
		UpdateExpression: aws.String("SET articleCursor = :articleCursor, lastBroadcastTimeMs = :lastBroadcastTimeMs"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":articleCursor": {
				N: aws.String(strconv.Itoa(subscriber.ArticleCursor)),
			},
			":lastBroadcastTimeMs": {
				N: aws.String(strconv.FormatInt(subscriber.LastBroadcastTimeMs, 10)),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdateSubscriberStat(): error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) DeleteSubscriber(projectId, email string) error {
	log.Infof("DeleteSubscriber(): projectId = %s, email = %s", projectId, email)
	deleteItemInput := dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"projectId": {S: aws.String(projectId)},
			"email":     {S: aws.String(email)},
		},
		TableName: d.subscribersTable,
	}
	_, err := d.dynamodbSvc.DeleteItem(&deleteItemInput)
	if err != nil {
		return err
	}
	return nil
}

func (d *DynamodbAdapter) ScanSubscriber(limit int, exclusiveStartKey map[string]*dynamodb.AttributeValue) (
	[]*models.Subscriber, map[string]*dynamodb.AttributeValue, error) {
	log.Infof("ScanSubscriber(): limit = %d, exclusiveStartKey = %v", limit, exclusiveStartKey)
	scanItemInput := dynamodb.ScanInput{
		Limit:     aws.Int64(int64(limit)),
		TableName: d.subscribersTable,
	}
	if len(exclusiveStartKey) > 0 {
		scanItemInput.ExclusiveStartKey = exclusiveStartKey
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanItemInput)
	if err != nil {
		log.Errorf("ScanSubscriber(): error = %v", err)
	}

	var res []*models.Subscriber
	for _, item := range dbOutput.Items {
		subscriber := &models.Subscriber{}
		err = dynamodbattribute.UnmarshalMap(item, subscriber)
		if err != nil {
			log.Warnf("ScanSubscriber(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, subscriber)
	}
	return res, dbOutput.LastEvaluatedKey, nil
}

func (d *DynamodbAdapter) ReadSubscriber(email, projectId string) (*models.Subscriber, error) {
	log.Infof("ReadSubscriber(): email = %s, projectId = %s", email, projectId)
	getItemInput := dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":     {S: aws.String(email)},
			"projectId": {S: aws.String(projectId)},
		},
		TableName: d.subscribersTable,
	}
	dbOutput, err := d.dynamodbSvc.GetItem(&getItemInput)
	if err != nil {
		log.Errorf("ReadSubscriber(): failed to get item, error: %+v", err)
		return nil, err
	}
	if dbOutput == nil || dbOutput.Item == nil {
		err = fmt.Errorf("item(%s, %s) not found", email, projectId)
		return nil, err
	}
	subscriber := &models.Subscriber{}
	err = dynamodbattribute.UnmarshalMap(dbOutput.Item, subscriber)
	if err != nil {
		log.Errorf("ReadSubscriber(): failed to unmarshall, error: %+v", err)
		return nil, err
	}
	return subscriber, nil
}

func (d *DynamodbAdapter) ReadSubscribers(projectId, token string, pageSize int) ([]*models.Subscriber, string, error) {
	log.Infof("ReadSubscribers(): projectId = %s, token = %s, pageSize = %d, table = %s", projectId, token, pageSize, *d.subscribersTable)
	queryInput := dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"projectId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(projectId),
					},
				},
			},
		},
		TableName: d.subscribersTable,
		Limit:     aws.Int64(int64(pageSize)),
	}
	if token != "" && strings.Contains(token, "#") {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"projectId": {
				S: aws.String(strings.Split(token, "#")[0]),
			},
			"email": {
				S: aws.String(strings.Split(token, "#")[1]),
			},
		}
	}

	dbOutput, err := d.dynamodbSvc.Query(&queryInput)
	if err != nil {
		log.Errorf("ReadSubscribers(): error = %v", err)
	}

	res := []*models.Subscriber{}
	for _, item := range dbOutput.Items {
		subscriber := &models.Subscriber{}
		err = dynamodbattribute.UnmarshalMap(item, subscriber)
		if err != nil {
			log.Warnf("ReadSubscribers(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, subscriber)
	}
	nextToken := "EOF"
	if len(dbOutput.LastEvaluatedKey) > 0 {
		nextToken = *dbOutput.LastEvaluatedKey["projectId"].S + "#" + *dbOutput.LastEvaluatedKey["email"].S
	}
	return res, nextToken, nil
}

func (d *DynamodbAdapter) SearchSubscribers(projectId, emailFilter string) ([]*models.Subscriber, error) {
	log.Infof("SearchSubscribers(): projectId = %s, emailFilter = %s", projectId, emailFilter)
	scanInput := dynamodb.ScanInput{
		ScanFilter: map[string]*dynamodb.Condition{
			"projectId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(projectId),
					},
				},
			},
			"email": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(emailFilter),
					},
				},
			},
		},
		Limit:     aws.Int64(100),
		TableName: d.subscribersTable,
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanInput)
	if err != nil {
		log.Errorf("SearchSubscribers(): error = %v", err)
	}

	res := []*models.Subscriber{}
	for _, item := range dbOutput.Items {
		subscriber := &models.Subscriber{}
		err = dynamodbattribute.UnmarshalMap(item, subscriber)
		if err != nil {
			log.Warnf("SearchSubscribers(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, subscriber)
	}

	return res, nil
}

func (d *DynamodbAdapter) UpdateUserSubscriptionPlan(email string, targetPlan, paymentStatus int,
	emailUsageInCycle int64, subscriptionId, subscriptionPriceId string, isBlock bool) error {
	log.Infof("UpdateUserSubscriptionPlan(): email = %s, targetPlan = %d, paymentStatus = %d, "+
		"emailUsageInCycle = %d, subscriptionId = %s, subscriptionPriceId = %s, isBlock = %s", email, targetPlan, paymentStatus,
		emailUsageInCycle, subscriptionId, subscriptionPriceId, isBlock)
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: d.usersTable,
		UpdateExpression: aws.String("SET subscriptionPlan = :targetPlan, paymentStatus = :paymentStatus, " +
			"emailUsageInCycle = :emailUsageInCycle, subscriptionId = :subscriptionId, " +
			"subscriptionPriceId = :subscriptionPriceId, isBlock = :isBlock"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":targetPlan": {
				N: aws.String(strconv.Itoa(targetPlan)),
			},
			":paymentStatus": {
				N: aws.String(strconv.Itoa(paymentStatus)),
			},
			":emailUsageInCycle": {
				N: aws.String(strconv.FormatInt(emailUsageInCycle, 10)),
			},
			":subscriptionId": {
				S: aws.String(subscriptionId),
			},
			":subscriptionPriceId": {
				S: aws.String(subscriptionPriceId),
			},
			"isBlock": {
				BOOL: aws.Bool(isBlock),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdateUserSubscriptionPlan(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) AddEmailUsageInCycle(email string, count int64) error {
	log.Infof("AddEmailUsageInCycle(): email = %s, count = %d", email, count)
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName:        d.usersTable,
		UpdateExpression: aws.String("SET emailUsageInCycle = emailUsageInCycle + :count"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":count": {
				N: aws.String(strconv.FormatInt(count, 10)),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("AddEmailUsageInCycle(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) UpdateAccountStatus(email string, isBlock bool, paymentStatus int) error {
	log.Infof("UpdateAccountStatus(): email = %s, isBlock = %s, paymentStatus = %d", email, isBlock, paymentStatus)
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName:        d.usersTable,
		UpdateExpression: aws.String("SET isBlock = :isBlock, paymentStatus = :paymentStatus"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":paymentStatus": {
				N: aws.String(strconv.Itoa(paymentStatus)),
			},
			":isBlock": {
				BOOL: aws.Bool(isBlock),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdateAccountStatus(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) UpdatePaymentStatus(email string, paymentStatus int) error {
	log.Infof("UpdatePaymentStatus(): email = %s, paymentStatus = %d", email, paymentStatus)
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName:        d.usersTable,
		UpdateExpression: aws.String("SET paymentStatus = :paymentStatus"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":paymentStatus": {
				N: aws.String(strconv.Itoa(paymentStatus)),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdatePaymentStatus(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) ChangePassword(email string, passwordHash string) error {
	log.Infof("ChangePassword(): email = %s", email)
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName:        d.usersTable,
		UpdateExpression: aws.String("SET passwordHash = :passwordHash"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":passwordHash": {
				S: aws.String(passwordHash),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("ChangePassword(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) ResetPassword(email string, timestamp int64) (string, error) {
	log.Infof("ResetPassword(): email = %s, timestamp = %d", email, timestamp)
	token, err := generateHash(fmt.Sprintf("%s#%d", email, timestamp))
	if err != nil {
		return "", err
	}
	resetPassToken := &models.ResetPasswordToken{
		Token:     token,
		Email:     email,
		CreatedAt: timestamp,
		ExpiredAt: time.Now().AddDate(0, 0, 1).Unix(),
	}

	userIdExistsCondition := dynamodb.ConditionCheck{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
		},
		ConditionExpression: aws.String("attribute_exists(email)"),
		TableName:           d.usersTable,
	}

	tableItem, err := MarshalMap(resetPassToken)
	putItemInput := dynamodb.Put{
		Item:      tableItem,
		TableName: d.resetPasswordTable,
	}

	transactionItem1 := dynamodb.TransactWriteItem{
		ConditionCheck: &userIdExistsCondition,
	}
	transactionItem2 := dynamodb.TransactWriteItem{
		Put: &putItemInput,
	}

	transactionInput := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{&transactionItem1, &transactionItem2},
	}

	_, err = d.dynamodbSvc.TransactWriteItems(&transactionInput)
	if err != nil {
		log.Errorf("ResetPassword(): error: %+v", err)
		return "", err
	}

	return token, nil
}

func (d *DynamodbAdapter) ReadPasswordReset(token string) (*models.ResetPasswordToken, error) {
	log.Infof("ReadPasswordReset(): token = %s", token)
	getItemInput := dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"token": {S: aws.String(token)},
		},
		TableName: d.resetPasswordTable,
	}
	dbOutput, err := d.dynamodbSvc.GetItem(&getItemInput)
	if err != nil {
		log.Errorf("ReadPasswordReset(): failed to get item, error: %+v", err)
		return nil, err
	}
	if dbOutput == nil || dbOutput.Item == nil {
		err = fmt.Errorf("item(%s) not found", token)
		return nil, err
	}
	resetToken := &models.ResetPasswordToken{}
	err = dynamodbattribute.UnmarshalMap(dbOutput.Item, resetToken)
	if err != nil {
		log.Errorf("ReadUser(): failed to unmarshall, error: %+v", err)
		return nil, err
	}
	return resetToken, nil
}

func (d *DynamodbAdapter) UpdatePassword(email, passwordHash string) error {
	log.Infof("UpdatePassword():")
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName:        d.usersTable,
		UpdateExpression: aws.String("SET passwordHash = :passwordHash"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":passwordHash": {
				S: aws.String(passwordHash),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdatePassword(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) CreateTicket(ticket *models.Ticket) error {
	log.Infof("CreateTicket(): ticketId = %s", ticket.TicketId)
	tableItem, err := MarshalMap(ticket)
	if err != nil {
		log.Errorf("CreateTicket(): failed to marshal user")
	}
	putItemInput := dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(ticketId)"),
		Item:                tableItem,
		TableName:           d.ticketTable,
	}
	_, err = d.dynamodbSvc.PutItem(&putItemInput)
	if err != nil {
		log.Errorf("CreateTicket(): failed to put item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) UpdateTicket(ticket *models.Ticket) error {
	log.Infof("UpdateTicket(): ticketId = %s", ticket.TicketId)
	tableItem, err := MarshalMap(ticket)
	if err != nil {
		log.Errorf("UpdateTicket(): failed to marshal user")
	}
	putItemInput := dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_exists(ticketId)"),
		Item:                tableItem,
		TableName:           d.ticketTable,
	}
	_, err = d.dynamodbSvc.PutItem(&putItemInput)
	if err != nil {
		log.Errorf("UpdateTicket(): failed to put item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) ReadTickets(email, token string, pageSize int) ([]*models.Ticket, string, error) {
	log.Infof("ReadTickets(): email = %s, token = %s, pageSize = %d", email, token, pageSize)
	queryInput := dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
		},
		TableName: d.ticketTable,
		Limit:     aws.Int64(int64(pageSize)),
	}
	if token != "" && strings.Contains(token, "#") {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"ticketId": {
				S: aws.String(strings.Split(token, "#")[0]),
			},
			"email": {
				S: aws.String(strings.Split(token, "#")[1]),
			},
		}
	}

	dbOutput, err := d.dynamodbSvc.Query(&queryInput)
	if err != nil {
		log.Errorf("ReadTickets(): error = %v", err)
	}

	res := []*models.Ticket{}
	for _, item := range dbOutput.Items {
		ticket := &models.Ticket{}
		err = dynamodbattribute.UnmarshalMap(item, ticket)
		if err != nil {
			log.Warnf("ReadTickets(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, ticket)
	}
	nextToken := "EOF"
	if len(dbOutput.LastEvaluatedKey) > 0 {
		nextToken = *dbOutput.LastEvaluatedKey["ticketId"].S + "#" + *dbOutput.LastEvaluatedKey["email"].S
	}
	return res, nextToken, nil
}

func (d *DynamodbAdapter) ScanTickets(limit int, token string) ([]*models.Ticket, string, error) {
	log.Infof("ScanTickets(): limit = %d, token = %s", limit, token)
	scanInput := dynamodb.ScanInput{
		FilterExpression: aws.String("ticketStatus = :open"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":open": {
				N: aws.String(strconv.FormatInt(1, 10)), // open ticket only
			},
		},
		Limit:     aws.Int64(int64(limit)),
		TableName: d.ticketTable,
	}

	var exclusiveStartKey map[string]*dynamodb.AttributeValue
	if token != "" && strings.Contains(token, "#") {
		exclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"ticketId": {
				S: aws.String(strings.Split(token, "#")[0]),
			},
			"email": {
				S: aws.String(strings.Split(token, "#")[1]),
			},
		}
		scanInput.ExclusiveStartKey = exclusiveStartKey
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanInput)
	if err != nil {
		log.Errorf("ScanTickets(): error = %v", err)
	}

	var res []*models.Ticket
	for _, item := range dbOutput.Items {
		ticket := &models.Ticket{}
		err = dynamodbattribute.UnmarshalMap(item, ticket)
		if err != nil {
			log.Warnf("ScanTickets(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, ticket)
	}

	nextToken := ""
	if len(dbOutput.LastEvaluatedKey) > 0 {
		nextToken = dbOutput.LastEvaluatedKey["email"].String() + "#" + dbOutput.LastEvaluatedKey["ticketId"].String()
	} else {
		nextToken = "EOF"
	}

	return res, nextToken, nil
}

func (d *DynamodbAdapter) DeleteTicket(email, ticketId string) error {
	log.Infof("DeleteTicket(): email = %s, ticketId = %s", email, ticketId)
	deleteItemInput := dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ticketId": {S: aws.String(ticketId)},
			"email":    {S: aws.String(email)},
		},
		TableName: d.ticketTable,
	}
	_, err := d.dynamodbSvc.DeleteItem(&deleteItemInput)
	if err != nil {
		return err
	}
	return nil
}

func (d *DynamodbAdapter) AppendComment(email, ticketId string, comment *models.Comment) error {
	log.Infof("AppendComment(): ticketId = %s, commentId = %s", ticketId, comment.CommendId)
	tableItem, err := MarshalMap(comment)
	if err != nil {
		log.Errorf("AppendComment(): failed to marshal user")
	}
	updateItemInput := dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"ticketId": {
				S: aws.String(ticketId),
			},
		},
		ConditionExpression: aws.String("attribute_exists(email) AND attribute_exists(ticketId)"),
		UpdateExpression:    aws.String("SET comments = list_append(if_not_exists(comments, :empty_list), :comment)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":comment": {
				L: []*dynamodb.AttributeValue{
					{
						M: tableItem,
					},
				},
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: d.ticketTable,
	}
	_, err = d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("AppendComment(): failed to update item, error: %+v", err)
		return err
	}
	return nil
}

func (d *DynamodbAdapter) SearchUsers(emailFilter string) ([]*models.User, error) {
	log.Infof("SearchUsers(): emailFilter = %s", emailFilter)
	scanInput := dynamodb.ScanInput{
		ScanFilter: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(emailFilter),
					},
				},
			},
		},
		Limit:     aws.Int64(100),
		TableName: d.usersTable,
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanInput)
	if err != nil {
		log.Errorf("SearchUsers(): error = %v", err)
	}

	res := []*models.User{}
	for _, item := range dbOutput.Items {
		user := &models.User{}
		err = dynamodbattribute.UnmarshalMap(item, user)
		if err != nil {
			log.Warnf("SearchUsers(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, user)
	}

	return res, nil
}

func (d *DynamodbAdapter) SearchTickets(ticketFilter string) ([]*models.Ticket, error) {
	log.Infof("SearchTickets(): ticketFilter = %s", ticketFilter)
	scanInput := dynamodb.ScanInput{
		ConditionalOperator: aws.String("OR"),
		ScanFilter: map[string]*dynamodb.Condition{
			"title": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ticketFilter),
					},
				},
			},
			"projectId": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ticketFilter),
					},
				},
			},
			"ticketId": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ticketFilter),
					},
				},
			},
			"projectName": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ticketFilter),
					},
				},
			},
			"email": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ticketFilter),
					},
				},
			},
			"name": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ticketFilter),
					},
				},
			},
			"body": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ticketFilter),
					},
				},
			},
		},
		Limit:     aws.Int64(100),
		TableName: d.ticketTable,
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanInput)
	if err != nil {
		log.Errorf("SearchTickets(): error = %v", err)
	}

	res := []*models.Ticket{}
	for _, item := range dbOutput.Items {
		ticket := &models.Ticket{}
		err = dynamodbattribute.UnmarshalMap(item, ticket)
		if err != nil {
			log.Warnf("SearchTickets(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, ticket)
	}

	return res, nil
}

// Admin updates user account feature
func (d *DynamodbAdapter) UpdateEditedUserInfo(editedUserInfo *models.EditedUserInfo) (*models.UserAccountInfo, error) {
	log.Infof("UpdateEditedUserInfo():")
	updateItemInput := dynamodb.UpdateItemInput{
		ConditionExpression: aws.String("attribute_exists(email)"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(editedUserInfo.Email),
			},
		},
		TableName: d.usersTable,
		UpdateExpression: aws.String("SET isBlock = :isBlock, userScope = :userScope, subscriptionPlan = :subscriptionPlan," +
			" emailUsageInCycle = :emailUsageInCycle, stripeCustomerId = :stripeCustomerId, paymentStatus = :paymentStatus"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isBlock": {
				BOOL: aws.Bool(editedUserInfo.IsBlock),
			},
			":userScope": {
				N: aws.String(strconv.FormatInt(int64(editedUserInfo.UserScope), 10)),
			},
			":subscriptionPlan": {
				N: aws.String(strconv.FormatInt(int64(editedUserInfo.SubscriptionPlan), 10)),
			},
			":emailUsageInCycle": {
				N: aws.String(strconv.FormatInt(int64(editedUserInfo.EmailUsageInCycle), 10)),
			},
			":stripeCustomerId": {
				S: aws.String(editedUserInfo.StripeCustomerId),
			},
			":paymentStatus": {
				N: aws.String(strconv.FormatInt(int64(editedUserInfo.PaymentStatus), 10)),
			},
		},
	}
	_, err := d.dynamodbSvc.UpdateItem(&updateItemInput)
	if err != nil {
		log.Errorf("UpdateEditedUserInfo(): failed to update item, error: %+v", err)
		return nil, err
	}

	updatedUser, err := d.ReadUser(editedUserInfo.Email)
	if err != nil {
		return nil, err
	}
	userAccountInfo := models.NewUserAccountInfo(updatedUser)

	return userAccountInfo, nil
}

func (d *DynamodbAdapter) SearchSubscribersByKey(subscriberFilter string) ([]*models.Subscriber, error) {
	log.Infof("SearchSubscribersByKey(): subscriberFilter = %s", subscriberFilter)
	scanInput := dynamodb.ScanInput{
		ConditionalOperator: aws.String("OR"),
		ScanFilter: map[string]*dynamodb.Condition{
			"projectId": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(subscriberFilter),
					},
				},
			},
			"email": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(subscriberFilter),
					},
				},
			},
			"firstName": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(subscriberFilter),
					},
				},
			},
			"lastName": {
				ComparisonOperator: aws.String("CONTAINS"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(subscriberFilter),
					},
				},
			},
		},
		Limit:     aws.Int64(100),
		TableName: d.subscribersTable,
	}

	dbOutput, err := d.dynamodbSvc.Scan(&scanInput)
	if err != nil {
		log.Errorf("SearchSubscribersByKey(): error = %v", err)
	}

	res := []*models.Subscriber{}
	for _, item := range dbOutput.Items {
		subscriber := &models.Subscriber{}
		err = dynamodbattribute.UnmarshalMap(item, subscriber)
		if err != nil {
			log.Warnf("SearchSubscribersByKey(): failed to unmarshall, error: %+v", err)
			continue
		}
		res = append(res, subscriber)
	}

	return res, nil
}

func generateHash(original string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(original), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("generateHash(): error = %v", err)
		return "", err
	}
	return string(bytes), err
}

func isListChanged(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return true
	}
	memo := map[string]bool{}
	for _, item := range l1 {
		memo[item] = true
	}
	for _, item := range l2 {
		if _, ok := memo[item]; !ok {
			return true
		}
	}
	return false
}
