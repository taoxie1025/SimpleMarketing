<template>
    <v-card color="blue-grey darken-1" dark flat class="mx-auto">
      <v-card-title class="headline grey lighten-2">
        <v-icon v-if="projectCopy.createdAt > 0" left color="black">mdi-cog-outline</v-icon>
        <h3 v-if="projectCopy.createdAt > 0">Settings</h3>
        <h3 v-else>New</h3>
      </v-card-title>
        <v-form ref="DetailsForm" v-model="isValid" lazy-validation>
            <v-layout wrap>
                <v-flex xs12>
                    <v-progress-linear :active="isUpdating" class="ma-0" color="green lighten-3" height="4" indeterminate></v-progress-linear>
                </v-flex>
            </v-layout>
            <v-container>
                <v-layout wrap>
                    <v-flex xs12 md12>
                      <br>
                      <v-text-field v-model="projectCopy.name" :rules="projectNameRules" :disabled="isUpdating" color="blue-grey lighten-2" label="Project Name" dense></v-text-field>
                    </v-flex>
                    <v-flex xs12 md12>
                        <v-textarea v-model="projectCopy.intro" :disabled="isUpdating" color="blue-grey lighten-2" dense>
                            <template v-slot:label>
                                <div>
                                    Intro <small>(optional)</small>
                                </div>
                            </template>
                        </v-textarea>
                    </v-flex>
                  <v-flex xs12 md12>
                    <v-text-field v-model="projectCopy.author" :disabled="isUpdating" color="blue-grey lighten-2" label="Author" dense></v-text-field>
                  </v-flex>
                </v-layout>
            </v-container>
            <v-container>
                <v-layout wrap>
                    <v-flex xs12 md12>
                        <v-text-field v-model="projectCopy.outgoingEmail" :rules="loginEmailRules" :disabled="isUpdating" color="blue-grey lighten-2" label="Outgoing Email" dense></v-text-field>
                    </v-flex>
                    <v-flex xs12 md12>
                      <div class="input-button-wrap">
                        <v-text-field v-model="projectCopy.avatarUrl" :disabled="isUpdating" color="blue-grey lighten-2" dense>
                            <template v-slot:label>
                                <div>
                                    Avatar URL <small>(optional)</small>
                                </div>
                            </template>
                          <template slot="append">
                            <v-btn tile @click="uploadFile('avatar')" color="blue-grey darken-3" class="pa-2">
                              <v-icon left>mdi-cloud-outline</v-icon>
                              Upload
                            </v-btn>
                          </template>
                        </v-text-field>
                      </div>
                    </v-flex>
                    <v-flex xs12 md12>
                      <div class="input-button-wrap">
                        <v-text-field v-model="projectCopy.backgroundImageUrl" :disabled="isUpdating" color="blue-grey lighten-2" dense>
                            <template v-slot:label>
                                <div>
                                    Background URL <small>(optional)</small>
                                </div>
                            </template>
                            <template slot="append">
                              <v-btn tile @click="uploadFile('background')" color="blue-grey darken-3" class="pa-2">
                                <v-icon left>mdi-cloud-outline</v-icon>
                                Upload
                              </v-btn>
                            </template>
                        </v-text-field>
                      </div>
                    </v-flex>
                </v-layout>
            </v-container>
            <v-container>
                <v-row align="center">
                    <v-col class="d-flex" cols="12" sm="6">
                        <v-select dense outlined v-model="projectCopy.interval" :items="intervals" item-value="interval" item-text="label" label="Frequency"></v-select>
                    </v-col>
                  <v-col class="d-flex" cols="12" sm="6">
                    <v-select
                        v-model="projectCopy.subscriptionType"
                        placeholder="Default"
                        :hint="getSubscriptionTypeHint()"
                        :items="subscriptionTypes"
                        item-text="option"
                        item-value="value"
                        label="Subscription Type"
                        persistent-hint
                        outlined
                        dense
                    ></v-select>
                  </v-col>
                </v-row>
            </v-container>
            <v-container v-if="projectCopy.createdAt">
                <v-layout wrap>
                    <v-btn block tile :disabled="isUpdating" :loading="isUpdating" color="red" depressed @click="dialog=true">
                        <v-icon left>mdi-trash-can-outline</v-icon>
                        Delete
                    </v-btn>
                    <v-row justify="center">
                        <v-dialog v-model="dialog" max-width="250">
                            <confirm-dialog v-bind:title="`Are you sure?`" :body="`The action is not reversible. All of your articles and subscribers will be deleted.`" @yes="deleteConfirmed" @no="dialog=false"></confirm-dialog>
                        </v-dialog>
                    </v-row>
                </v-layout>
            </v-container>
            <v-container>
                <v-row justify="center">
                    <v-dialog v-model="sendVerificationDialog" max-width="650">
                        <confirm-dialog v-bind:title="compileVerificationDialog.title" :body="compileVerificationDialog.body" :simpleConfirm="true" @no="sendVerificationDialog=false"></confirm-dialog>
                    </v-dialog>
                </v-row>
            </v-container>
            <v-card-actions>
                <v-switch v-if="projectCopy.createdAt > 0" v-model="projectState" :disabled="isUpdating" class="mt-0" color="green lighten-2" hide-details flat :label="`${projectState ? 'Online' : 'Offline'}`"></v-switch>
                <v-spacer></v-spacer>
                <v-btn tile v-if="projectCopy.createdAt > 0" :disabled="isUpdating||!isValid" :loading="isUpdating" color="primary" @click="updateProject()">
                    <v-icon left>mdi-update</v-icon>
                    Save
                </v-btn>
                <v-btn v-else tile :disabled="isUpdating||!isValid" :loading="isUpdating" color="primary" @click="createProject()">
                    <v-icon left>mdi-update</v-icon>
                    Create
                </v-btn>
            </v-card-actions>
        </v-form>
        <v-container>
            <v-dialog v-model="showGallery" width="500px">
              <v-card>
                <file-selector @select-file="onSelectImage"></file-selector>
              </v-card>
            </v-dialog>
        </v-container>
    </v-card>
</template>

<script>
    import axios from "axios";
    import {default as API_ENDPOINTS} from "../api";
    import ConfirmDialog from "./ConfirmDialog";
    import FileSelector from './Gallery'

    export default {
        components: {ConfirmDialog, FileSelector},
        data() {
            return {
                projectCopy: {},
                projectState: '',
                isUpdating: false,
                dialog: false,
                sendVerificationDialog: false,
                loginEmailRules: [
                    v => !!v || "Required",
                    v => /.+@.+\..+/.test(v) || "E-mail must be valid"
                ],
                projectNameRules: [
                    v => !!v || "Required",
                ],
                isValid: false,
                intervals: [
                    {
                        interval: 1 * 24 * 60 * 60 * 1000,
                        label: "1 day"
                    },
                    {
                        interval: 2 * 24 * 60 * 60 * 1000,
                        label: "2 days"
                    },
                    {
                        interval: 3 * 24 * 60 * 60 * 1000,
                        label: "3 days"
                    },
                    {
                        interval: 4 * 24 * 60 * 60 * 1000,
                        label: "4 days"
                    },
                    {
                        interval: 5 * 24 * 60 * 60 * 1000,
                        label: "5 days"
                    },
                    {
                        interval: 6 * 24 * 60 * 60 * 1000,
                        label: "6 days"
                    },
                    {
                        interval: 7 * 24 * 60 * 60 * 1000,
                        label: "1 week"
                    },
                    {
                        interval: 14 * 24 * 60 * 60 * 1000,
                        label: "2 weeks"
                    },
                    {
                        interval: 21 * 24 * 60 * 60 * 1000,
                        label: "3 weeks"
                    },
                    {
                        interval: 30 * 24 * 60 * 60 * 1000,
                        label: "1 month"
                    },
                    {
                        interval: 60 * 24 * 60 * 60 * 1000,
                        label: "2 months"
                    },
                ],
                showGallery: false,
                isEditAvatarUrl: false,
                isEditBackgroundUrl: false,
                subscriptionTypes: [
                  {option: "Default", value: 0},
                  {option: "Rolling", value: 1}
                ],
            }
        },
        watch: {
            projectState: function() {
                if (this.projectState) {
                    this.project.projectState = 3 // Live
                } else {
                    this.project.projectState = 2 // Created
                }
            },
            resetFromValidation: function () {
                this.$refs.DetailsForm.resetValidation()
            }
        },
        props:[
            'project',
            'resetFromValidation'
        ],
        computed: {
            compileVerificationDialog() {
                const dialog = {
                    title: "The email address " + this.projectCopy.outgoingEmail + " is not yet verified",
                    body: "A verification email has been sent to the email, please go to your mailbox and verify the address before proceeding."
                }
                return dialog
            },
            validateForm: function() {
                if (this.$refs.DetailsForm.validate()) {
                    return true
                }
                return false
            },
            lastBroadcastTime: function () {
                if (this.project.lastBroadcastTimeMs && this.project.lastBroadcastTimeMs > 0) {
                    const date = new Date(this.project.lastBroadcastTimeMs)
                    return date.toLocaleDateString()
                }
                return ""
            },
            createdTime: function () {
                if (this.project.createdAt && this.project.createdAt > 0) {
                    const date = new Date(this.project.createdAt)
                    return date.toLocaleDateString()
                }
                return ""
            }
        },
        methods: {
            getSubscriptionTypeHint() {
              if (this.projectCopy.subscriptionType == 0) {
                return "Start from the first article"
              } else {
                return "Start from the next new article"
              }
            },
            createProject() {
                if (!this.validateForm) {
                   return
                }
                this.isOutgoingEmailVerified().then((resp) => {
                    if (resp?.data?.isVerified) {
                        this.isUpdating = true
                        return new Promise((resolve, reject) => {
                            axios({url: API_ENDPOINTS.CREATE_PROJECT, data: this.project, method: 'POST'})
                                .then(resp => {
                                    resolve(resp)
                                    this.isUpdating = false
                                    this.$emit('projectCreated', (resp.data))
                                    this.$refs.DetailsForm.resetValidation()
                                })
                                .catch(err => {
                                    reject(err)
                                    this.isUpdating = false
                                })
                        })
                    } else {
                        this.sendVerificationDialog= true
                        this.sendVerificationEmail(this.projectCopy.outgoingEmail)
                    }
                }).catch(() => {

                })
            },
            updateProject() {
                if (!this.validateForm) {
                    return
                }

                this.isOutgoingEmailVerified().then((resp) => {
                    if (resp?.data?.isVerified) {
                        this.isUpdating = true
                        if (!this.project.articleIds || this.project.articleIds.length == 0) {
                          delete this.project["articleIds"]
                        }
                        const data = {email: this.$store.getters.email, project: this.project}
                        return new Promise((resolve, reject) => {
                            axios({url: API_ENDPOINTS.UPDATE_PROJECT, data: data, method: 'PUT' })
                                .then(resp => {
                                    resolve(resp)
                                    this.isUpdating = false
                                    this.$emit('projectUpdated', resp.data)
                                })
                                .catch(err => {
                                    reject(err)
                                    this.isUpdating = false
                                })
                        })
                    } else {
                        this.sendVerificationDialog = true
                        this.sendVerificationEmail(this.projectCopy.outgoingEmail)
                    }
                }).catch(() => {

                })
            },
            isOutgoingEmailVerified() {
                const data = {requesterEmail: this.$store.getters.email}
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.IS_EMAIL_VERIFIED(this.projectCopy.outgoingEmail), params: data, method: 'GET' })
                        .then(resp => {
                            resolve(resp)
                        })
                        .catch((err) => {
                            reject(err)
                        })
                })
            },
            sendVerificationEmail(email) {
                const data = {requesterEmail: this.$store.getters.email}
                return new Promise((resolve) => {
                    axios({url: API_ENDPOINTS.SEND_VERIFICATION_EMAIL(email), data: data, method: 'POST' })
                        .then(resp => {
                            resolve(resp)
                        })
                        .catch(() => {
                        })
                })
            },
            deleteConfirmed() {
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.DELETE_PROJECT(this.project.email, this.project.projectId), method: 'DELETE' })
                        .then(resp => {
                            resolve(resp)
                            this.$emit('projectDeleted')
                            this.dialog = false
                        })
                        .catch(err => {
                            reject(err)
                            this.dialog = false
                        })
                })
            },
          uploadFile(field) {
            this.showGallery = true
            if (field == 'avatar') {
              this.isEditAvatarUrl = true
            } else if (field == 'background') {
              this.isEditBackgroundUrl = true
            }
          },
          onSelectImage(image) {
            if (this.isEditAvatarUrl) {
              this.projectCopy.avatarUrl = image.src
            } else if (this.isEditBackgroundUrl) {
              this.projectCopy.backgroundImageUrl = image.src
            }
            this.isEditAvatarUrl = false
            this.isEditBackgroundUrl = false
            this.showGallery = false
          }
        },
        beforeUpdate() {
            this.projectCopy = this.project
            if (this.projectCopy && (!this.projectCopy.outgoingEmail || this.projectCopy.outgoingEmail == "")) {
                if (this.projectCopy.email) {
                    this.projectCopy.outgoingEmail = this.projectCopy.email
                } else {
                    this.projectCopy.outgoingEmail = this.$store.getters.email
                }
            }
            if (!this.projectCopy.author && !this.projectCopy.createdAt) {
              this.projectCopy.author = this.$store.getters.user.firstName + " " + this.$store.getters.user.lastName
            }
            if (this.projectCopy && this.projectCopy.projectState == 3) {
                this.projectState = true
            } else {
                this.projectState = false
            }
        },
        created() {
            this.projectCopy = this.project
            if (this.projectCopy && (!this.projectCopy.outgoingEmail || this.projectCopy.outgoingEmail == "")) {
                if (this.projectCopy.email) {
                    this.projectCopy.outgoingEmail = this.projectCopy.email
                } else {
                    this.projectCopy.outgoingEmail = this.$store.getters.email
                }
            }
            if (this.projectCopy && this.projectCopy.projectState == 3) {
                this.projectState = true
            }
            if (!this.projectCopy.interval || this.projectCopy.interval < this.intervals[0].interval) {
                this.projectCopy.interval = this.intervals[0].interval
            }
        }
    }
</script>

<style lang="scss">
.input-button-wrap {
  display: flex;
}
</style>