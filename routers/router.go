package routers

import (
	"email_action/controllers"
	"email_action/mail"
	"email_action/store"
	"github.com/astaxie/beego"
)

func InitRouters(store *store.Store, mailer *mail.SesAdapter) {
	beego.Router("/", &controllers.HomePageController{Store: store})

	beego.AddNamespace(
		beego.NewNamespace("/api",
			beego.NSNamespace("/v1",
				beego.NSRouter("/authenticate", &controllers.AuthController{Store: store, Mailer: mailer}, "post:AuthUser"),
				beego.NSRouter("/authenticate/reset", &controllers.AuthController{Store: store, Mailer: mailer}, "post:ResetPassword"),
				beego.NSRouter("/authenticate/token", &controllers.AuthController{Store: store, Mailer: mailer}, "post:IsTokenValid"),
				beego.NSRouter("/authenticate/password", &controllers.AuthController{Store: store, Mailer: mailer}, "post:UpdatePassword"),
				beego.NSRouter("/register", &controllers.UserController{Store: store}, "post:CreateUser"),

				beego.NSRouter("/user/:email", &controllers.UserController{Store: store}, "put:UpdateUserBasicInfo"),
				beego.NSRouter("/user/change_password", &controllers.UserController{Store: store}, "put:UpdatePassword"),
				beego.NSRouter("/user/:email", &controllers.UserController{Store: store}, "get:ReadUser"),
				beego.NSRouter("/project", &controllers.ProjectController{Store: store}, "post:CreateProject"),
				beego.NSRouter("/project", &controllers.ProjectController{Store: store}, "put:UpdateProject"),
				beego.NSRouter("/projects", &controllers.ProjectController{Store: store}, "post:ReadProjects"),
				beego.NSRouter("/project/:email/:projectId", &controllers.ProjectController{Store: store}, "delete:DeleteProject"),
				beego.NSRouter("/project/:email/:projectId", &controllers.ProjectController{Store: store}, "get:ReadProject"),
				beego.NSRouter("/article", &controllers.ArticleController{Store: store}, "post:CreateArticle"),
				beego.NSRouter("/article", &controllers.ArticleController{Store: store}, "put:UpdateArticle"),
				beego.NSRouter("/article/:email/:projectId/:articleId", &controllers.ArticleController{Store: store}, "delete:DeleteArticle"),
				beego.NSRouter("/article/:email/:articleId", &controllers.ArticleController{Store: store}, "get:ReadArticle"),
				beego.NSRouter("/articles", &controllers.ArticleController{Store: store}, "Post:ReadArticles"),
				beego.NSRouter("/subscriber", &controllers.SubscriberController{Store: store}, "post:CreateSubscriber"),
				beego.NSRouter("/subscriber", &controllers.SubscriberController{Store: store}, "put:UpdateSubscriber"),
				beego.NSRouter("/subscriber/:projectId/:email", &controllers.SubscriberController{Store: store}, "delete:DeleteSubscriber"),
				beego.NSRouter("/subscribers", &controllers.SubscriberController{Store: store}, "post:ReadSubscribers"),
				beego.NSRouter("/subscribers", &controllers.SubscriberController{Store: store}, "get:SearchSubscribers"),

				beego.NSRouter("/file", &controllers.FileController{Store: store}, "Post:FileUpload"),
				beego.NSRouter("/file/:email", &controllers.FileController{Store: store}, "Get:ListFiles"),
				beego.NSRouter("/file/:email/:fileName", &controllers.FileController{Store: store}, "Delete:DeleteFile"),

				beego.NSRouter("/email/verify/:email", &controllers.EmailController{Store: store, Mailer: mailer}, "Post:SendVerificationEmail"),
				beego.NSRouter("/email/verify/:email", &controllers.EmailController{Store: store, Mailer: mailer}, "Get:IsEmailVerified"),
				beego.NSRouter("/email/contactus", &controllers.EmailController{Store: store, Mailer: mailer}, "Post:ContactUs"),
				beego.NSRouter("/email/send", &controllers.EmailController{Store: store, Mailer: mailer}, "Post:SendEmail"),

				beego.NSRouter("/brief/:projectId", &controllers.ProjectController{Store: store}, "Get:ReadProjectBrief"),

				beego.NSRouter("/checkout/session", &controllers.CheckoutController{Store: store}, "Post:CreateSession"),
				beego.NSRouter("/checkout/coupon", &controllers.CheckoutController{Store: store}, "Get:ReadCoupon"),
				beego.NSRouter("/checkout/:email/:subscriptionId", &controllers.CheckoutController{Store: store}, "Delete:CancelSubscription"),
				beego.NSRouter("/checkout/:email/:subscriptionId", &controllers.CheckoutController{Store: store}, "Put:UpdateSubscription"),
				beego.NSRouter("/webhook", &controllers.WebhookController{Store: store}, "Post:StripeEventListener"),

				beego.NSRouter("/account/:email", &controllers.UserController{Store: store}, "Get:GetAccountInfo"),

				beego.NSRouter("/ticket", &controllers.TicketController{Store: store, Mailer: mailer}, "post:CreateTicket"),
				beego.NSRouter("/ticket/:email/:ticketId", &controllers.TicketController{Store: store, Mailer: mailer}, "delete:DeleteTicket"),
				beego.NSRouter("/tickets", &controllers.TicketController{Store: store, Mailer: mailer}, "post:ReadTickets"),
				beego.NSRouter("/ticket/:ticketId", &controllers.TicketController{Store: store, Mailer: mailer}, "put:UpdateTicket"),
				beego.NSRouter("/comment", &controllers.TicketController{Store: store, Mailer: mailer}, "post:CreateComment"),

				beego.NSRouter("/admin/tickets", &controllers.AdminController{Store: store, Mailer: mailer}, "post:AdminReadTickets"),
				beego.NSRouter("/admin/tickets", &controllers.AdminController{Store: store}, "get:AdminSearchTickets"),
				beego.NSRouter("/admin/users", &controllers.AdminController{Store: store, Mailer: mailer}, "post:AdminReadUsers"),
				beego.NSRouter("/admin/users", &controllers.AdminController{Store: store}, "get:AdminSearchUsers"),
				beego.NSRouter("/admin/user", &controllers.AdminController{Store: store}, "post:AdminUpdateUser"),
				beego.NSRouter("/admin/subscribers", &controllers.AdminController{Store: store, Mailer: mailer}, "post:AdminReadSubscribers"),
				beego.NSRouter("/admin/subscribers", &controllers.AdminController{Store: store, Mailer: mailer}, "get:AdminSearchSubscribers"),
				beego.NSRouter("/admin/subscriber", &controllers.AdminController{Store: store, Mailer: mailer}, "put:AdminUpdateSubscriber"),
			),
		),
	)
}
