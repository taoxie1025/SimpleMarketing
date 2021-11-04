<template>

    <div id="app" >
        <v-app id="inspire" >
            <div v-if="isLoading">
                <v-row justify="center">
                    <v-progress-circular indeterminate color="primary"></v-progress-circular>
                </v-row>
            </div>
            <div v-else>
                <v-container class="my-2" grid-list-sm>
                    <v-layout>
                        <v-flex xs12>
                            <v-btn tile class="info" @click="articleEditorDialog=true">
                                <v-icon dark>mdi-plus</v-icon>
                                Add New Article
                            </v-btn>
                        </v-flex>
                    </v-layout>
                    <br>
                    <draggable v-model="articles" @change="saveOrder" handle=".handle" v-if="this.articles && this.articles.length > 0">
                        <v-card :disabled="isSaving" text v-for="(article, index) in articles" :key="article.articleId" class="handle">
                            <v-layout row wrap class="pa-2">
                                <v-flex xs2 sm2 md1>
                                    <div class="caption grey--text ml-9">No.</div>
                                    <div :class="{'mt-0 mb-1': $vuetify.breakpoint.smAndDown}">
                                      <v-tooltip bottom>
                                        <template v-slot:activator="{ on, attrs }">
                                          <v-icon size="35" color="primary" v-bind="attrs" v-on="on" >drag_indicator</v-icon>
                                        </template>
                                        <span>Drag to move</span>
                                      </v-tooltip>
                                        {{index + 1}}
                                    </div>
                                </v-flex>

                                <v-flex xs4 sm4 md4>
                                    <div class="caption grey--text mb-1">Title</div>
                                    <div :class="{'mt-0 mb-1': $vuetify.breakpoint.smAndDown}" id="title-container">
                                        {{ article.title }}
                                    </div>
                                </v-flex>

                                <v-flex xs2 sm2 md3>
                                    <div class="caption grey--text mb-1">Updated on</div>
                                    <div :class="{'mt-0 mx-0 px-0': $vuetify.breakpoint.smAndDown}" v-html="convertToDate(article.updatedAt)"></div>
                                </v-flex>

                                <v-flex xs1 sm2 md2 class="pt-0">
                                    <div class="mt-0 pt-0">
                                        <v-switch v-model="article.isLive" @click="changeState(index)" :label="`${article.isLive ? 'ON' : 'OFF'}`"></v-switch>
                                    </div>
                                </v-flex>
                                <v-spacer></v-spacer>
                                <v-flex xs2 sm2 md2 class="mt-0 pt-3" id="buttonContainer">
                                    <v-row class="ml-10">
                                    <div>
                                      <v-tooltip bottom>
                                        <template v-slot:activator="{ on, attrs }">
                                          <v-btn icon :class="article" @click="openArticleEditor(index)" v-bind="attrs" v-on="on" class="pa-3 ma-0"><v-icon size="35" color="primary">edit</v-icon></v-btn>
                                        </template>
                                        <span>Edit</span>
                                      </v-tooltip>
                                      <v-tooltip bottom>
                                        <template v-slot:activator="{ on, attrs }">
                                          <v-btn icon :class="article" @click="openSendTestEmailDialog(index)" v-bind="attrs" v-on="on" class="pa-3 ma-0"><v-icon size="35" color="green">mdi-email-send-outline</v-icon></v-btn>
                                        </template>
                                        <span>Send</span>
                                      </v-tooltip>
                                      <v-tooltip bottom>
                                        <template v-slot:activator="{ on, attrs }">
                                          <v-btn icon :class="article" @click="openDialog(index)" v-bind="attrs" v-on="on"  class="pa-0 ma-0"><v-icon size="35" color="grey">delete</v-icon></v-btn>
                                        </template>
                                        <span>Delete</span>
                                      </v-tooltip>
                                    </div>
                                    </v-row>
                                </v-flex>
                            </v-layout>
                        </v-card>
                    </draggable>
                    <v-row justify="center">
                        <v-dialog v-model="dialog" max-width="250">
                            <confirm-dialog v-bind:title="`Are you sure?`" :body="`The action is not reversible.`" @yes="deleteArticle" @no="dialog=false"></confirm-dialog>
                        </v-dialog>
                    </v-row>
                </v-container>
                <v-container>
                    <v-row justify="center">
                        <v-dialog v-if="articleEditorDialog" v-model="articleEditorDialog" fullscreen hide-overlay transition="dialog-bottom-transition">
                            <article-editor v-bind:article="articles && selectedIndex >= 0 ? Object.assign({}, articles[selectedIndex]) : {}" @save="onSave" @create="onCreate" @closeEditor="onCloseEditor"></article-editor>
                        </v-dialog>
                    </v-row>
                </v-container>
                <v-container>
                  <v-row justify="center">
                    <v-dialog v-if="sendTestEmailDialog" v-model="sendTestEmailDialog" max-width="350">
                      <email-address-card :email-address="project.outgoingEmail" @send="sendTestEmail" @close="closeSendTestEmailDialog"></email-address-card>
                    </v-dialog>
                  </v-row>
                </v-container>
                <v-container>
                  <v-row>
                    <v-snackbar v-model="messageSnackbar">
                      {{ displayMessage }}
                      <template v-slot:action="{ attrs }">
                        <v-btn
                            color="pink"
                            text
                            v-bind="attrs"
                            @click="snackbar = false"
                        >
                          Close
                        </v-btn>
                      </template>
                    </v-snackbar>
                  </v-row>
                </v-container>
            </div>
        </v-app>
    </div>

</template>

<script>

    import axios from "axios";
    import {default as API_ENDPOINTS} from "../api";
    import Draggable from 'vuedraggable';
    import ConfirmDialog from "../components/ConfirmDialog";
    import ArticleEditor from "../components/ArticleEditor";
    import EmailAddressCard from "@/components/EmailAddressCard";

    export default {
        components: {EmailAddressCard, Draggable, ConfirmDialog, ArticleEditor},
        data() {
            return {
                projectId: '',
                project: {},
                articles: [],
                dialog: false,
                selectedIndex: -1,
                isSaving: false,
                articleEditorDialog: false,
                isLoading: false,
                sendTestEmailDialog: false,
                messageSnackbar: false,
                displayMessage: ""
            }
        },
        watch: {
        },
        methods: {
            openSendTestEmailDialog(index) {
              this.sendTestEmailDialog = true
              this.selectedIndex = index
            },
            closeSendTestEmailDialog() {
              this.sendTestEmailDialog = false
              this.selectedIndex = -1
            },
            sendTestEmail(to) {
              this.sendEmail(this.project.outgoingEmail, to, this.articles[this.selectedIndex].title,
                  this.articles[this.selectedIndex].htmlBody).then(() => {
                  this.closeSendTestEmailDialog()
                  this.showSnackBar("Message sent to " + to)
              }).catch(() => {
                  this.closeSendTestEmailDialog()
                  this.showSnackBar("Failed to send message to " + to)
              })
            },
            showSnackBar(msg) {
              this.displayMessage = msg
              this.messageSnackbar = true
              setTimeout(function(){
                this.messageSnackbar = false
              }, 2000);
            },
            onSave(editedArticle) {
                this.updateArticle(editedArticle).then(resp => {
                    this.articles.splice(this.selectedIndex, 1, resp.data)
                    this.articleEditorDialog = false
                    this.selectedIndex = -1
                })
            },
            onCreate(article) {
                this.createArticle(article).then(resp => {
                    this.articles.push(resp.data)
                    this.articleEditorDialog = false
                })
            },
            onCloseEditor() {
                this.articleEditorDialog = false
                this.selectedIndex = -1
            },
            changeState(index) {
                this.UpdateArticleByIndex(index)
            },
            UpdateArticleByIndex(index) {
                const article = this.articles[index]
                this.updateArticle(article).then(resp => {
                    this.articles.splice(index, 1, resp.data)
                })
            },
            createArticle(article) {
                const data = {email: this.project.email, projectId: this.projectId, isLive: article.isLive, title: article.title, htmlBody: article.htmlBody, textBody: article.textBody}
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.CREATE_ARTICLE, data: data, method: 'POST' })
                        .then(resp => {
                            resolve(resp)
                            return resp
                        })
                        .catch(err => {
                            reject(err)
                        })
                })
            },
            updateArticle(article) {
                const data = {email: this.$store.getters.email, article: article}
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.UPDATE_ARTICLE, data: data, method: 'PUT' })
                        .then(resp => {
                            resolve(resp)
                            return resp
                        })
                        .catch(err => {
                            reject(err)
                        })
                })
            },
            convertToDate(timestampMs) {
                if (timestampMs > 0) {
                    const date = new Date(timestampMs)
                    return date.toLocaleDateString()
                }
                return ""
            },
            saveOrder() {
                let newArticleIds = []
                for (let i = 0; i < this.articles.length; i++) {
                    newArticleIds.push(this.articles[i].articleId)
                }
                this.project.articleIds = newArticleIds
                this.saveProject()
            },
            saveProject() {
                this.isSaving = true
                const data = {email: this.$store.getters.email, project: this.project}
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.UPDATE_PROJECT, data: data, method: 'PUT' })
                        .then(resp => {
                            resolve(resp)
                            this.project = resp.data
                            this.readArticles()
                            this.isSaving = false
                        })
                        .catch(err => {
                            reject(err)
                        })
                })
            },
            openArticleEditor(index) {
                this.selectedIndex = index
                this.articleEditorDialog = true
            },
            openDialog(index) {
                this.dialog = true
                this.selectedIndex = index
            },
            deleteArticle() {
                const projectId = this.articles[this.selectedIndex].projectId
                const articleId = this.articles[this.selectedIndex].articleId
                const data = {requesterEmail: this.$store.getters.email}
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.DELETE_ARTICLE(this.$store.getters.email, projectId, articleId), params: data, method: 'DELETE' })
                        .then(resp => {
                            resolve(resp)
                            this.project.articleIds.splice(this.selectedIndex, 1)
                            this.articles.splice(this.selectedIndex, 1)
                            this.dialog = false
                        })
                        .catch(err => {
                            reject(err)
                        })
                })
            },
            readProject() {
                this.isLoading = true
                const data = {requesterEmail: this.$store.getters.email}
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.READ_PROJECT(this.$store.getters.email, this.projectId), params: data, method: 'GET' })
                        .then(resp => {
                            resolve(resp)
                            this.project = resp.data
                            this.readArticles()
                        })
                        .catch(err => {
                            reject(err)
                        })
                })
            },
            readArticles() {
                return new Promise((resolve, reject) => {
                    let data = {email: this.$store.getters.email ? this.$store.getters.email : this.project.email}
                    if (this.project?.articleIds?.length > 0) {
                        data.articleIds = this.project.articleIds
                    }
                    axios({url: API_ENDPOINTS.READ_ARTICLES, data: data, method: 'POST' })
                        .then(resp => {
                            resolve(resp)
                            this.articles = resp.data
                            this.isLoading = false
                        })
                        .catch(err => {
                            reject(err)
                        })
                })
            },
            sendEmail(from, to, subject, message) {
              return new Promise((resolve, reject) => {
                let data = {email: this.$store.getters.email ? this.$store.getters.email : this.project.email,
                  from: from, to: to, message: message, subject: subject, projectId: this.project.projectId,
                  projectName: this.project.name}
                axios({url: API_ENDPOINTS.SEND_EMAIL, data: data, method: 'POST' })
                    .then(resp => {
                      resolve(resp)
                    })
                    .catch(err => {
                      reject(err)
                    })
              })
            },
            isLogin() {
                if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
                    console.log("not signed in ");
                    this.$router.push('/')
                    return false
                }
                return true
            },
        },
        created: function () {
            if (!this.isLogin() || !this.$route.params || !this.$route.params.projectId) {
              this.$router.push('/dashboard')
            }
            this.projectId = this.$route.params.projectId
            if (!this.$route.params.project) {
                this.readProject()
            } else {
                this.project = this.$route.params.project
                this.readArticles()
            }
        }
    }
</script>


<style>
    .handle {
        cursor: move;
    }
    #buttonContainer {
        margin-left: auto;
    }
    #title-container {
      max-width: 300px;
    }
</style>
