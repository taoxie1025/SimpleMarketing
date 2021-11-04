<template>
    <div id="app" >
        <v-app id="inspire" >
            <br>
            <div v-if="isLoading">
                <v-row justify="center">
                    <v-progress-circular indeterminate color="primary"></v-progress-circular>
                </v-row>
            </div>
            <div v-else>
                <v-container fluid>
                    <v-layout row wrap class="ml-4 mr-4">
                        <v-flex xs12>
                            <v-row>
                                <v-btn tile dark color="info" @click="showProjectDetailsDialog">
                                  <v-icon left>mdi-plus</v-icon>
                                  Create New
                                </v-btn>
                                <v-spacer></v-spacer>
                                <v-tooltip top>
                                    <template v-slot:activator="{ on, attrs }">
                                    <v-chip  v-bind="attrs" v-on="on" label outlined @click="goToPlanPage()">
                                        Email Usage: {{$store.getters.user.emailUsageInCycle}} / {{getTotalQuota()}}
                                    </v-chip>
                                    </template>
                                    <span>Usage will be reset to 0 every 30 days</span>
                                </v-tooltip>
                            </v-row>
                        </v-flex>
                    </v-layout>
                </v-container>
                <v-container fluid>
                    <v-layout row wrap>
                        <v-flex xs12 sm12 md6 lg4 v-for="(project, index) in projects" v-bind:key="project.projectId">
                          <project-card outlined tile v-bind:project="project" :index="index" @showProjectSettings="onShowProjectSettings" @showProjectArticles="onShowProjectArticles(index)"></project-card>
                        </v-flex>
                    </v-layout>
                </v-container>
                <v-row justify="center">
                    <v-dialog v-model="NewProjectDialog" max-width="800px">
                        <project-details v-if="NewProjectDialog" v-bind:project="{email:this.$store.getters.email}" :resetFromValidation="resetFromValidation" @projectCreated="onProjectCreated"></project-details>
                    </v-dialog>
                </v-row>
                <v-row justify="center">
                    <v-dialog v-model="settingDialog" max-width="800px">
                        <project-details v-if="projects[selectedIndex] && settingDialog" v-bind:project="Object.assign({}, projects[selectedIndex])" :resetFromValidation="resetFromValidation" @projectDeleted="onProjectDeleted" @projectUpdated="onProjectUpdated"></project-details>
                    </v-dialog>
                </v-row>
            </div>
        </v-app>
    </div>

</template>

<script>

    import ProjectCard from "../components/ProjectCard";
    import ProjectDetails from "../components/ProjectDetails";
    import axios from "axios";
    import {default as API_ENDPOINTS} from "../api";
    export default {
        components: {ProjectCard, ProjectDetails},
        data() {
            return {
                projects: [],
                showProjectDetails: false,
                NewProjectDialog: false,
                settingDialog: false,
                selectedIndex: -1,
                resetFromValidation: false,
                isLoading: false
            }
        },
        watch: {
            NewProjectDialog(val) {
                if (!val) {
                    this.resetFromValidation = true
                } else {
                    this.resetFromValidation = false
                }
            },
            settingDialog(val) {
                if (!val) {
                    this.resetFromValidation = true
                } else {
                    this.resetFromValidation = false
                }
            }
        },
        methods: {
            goToPlanPage() {
              this.$router.push('/plan')
            },
            getTotalQuota() {
              switch (this.$store.getters.user.subscriptionPlan) {
                case 1: return "100000"
                case 2: return "300000"
              }
              return "1000"
            },
            readUserAndProjects() {
                this.$store.dispatch('readUser', this.$store.getters.email)
                    .then(() => {
                      if (this.$store.getters.user.projectIds && this.$store.getters.user.projectIds.length > 0) {
                          this.readProjects(this.$store.getters.email, this.$store.getters.user.projectIds).then(resp => {
                              const sortedProjects = resp.data.sort(function(x, y) {
                                  return x.createdAt < y.createdAt ? 1 : (x.createdAt > y.createdAt ? -1 : 0)
                              })
                              this.projects = sortedProjects
                              this.isLoading = false
                          })
                        } else {
                        this.isLoading = false
                      }
                    })
            },
            isLogin() {
                if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
                    this.$router.push('/')
                    return false
                }
                return true
            },
            readProjects(email, projectIds){
                const data = {email, projectIds}
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.READ_PROJECTS, data: data, method: 'POST' })
                        .then(resp => {
                            resolve(resp)
                        })
                        .catch(err => {
                            reject(err)
                        })
                })
            },
            showProjectDetailsDialog() {
                this.NewProjectDialog = true
            },
            onProjectCreated(project){
                this.projects.push(project)
                this.NewProjectDialog = false
                const sortedProjects = this.projects.sort(function(x, y) {
                    return x.createdAt < y.createdAt ? 1 : (x.createdAt > y.createdAt ? -1 : 0)
                })
                this.projects = sortedProjects
            },
            onProjectUpdated(project) {
                this.settingDialog = false
                this.projects.splice(this.selectedIndex, 1, project);
            },
            onProjectDeleted() {
                this.projects.splice(this.selectedIndex, 1);
                this.settingDialog = false
            },
            onShowProjectSettings(index) {
                this.selectedIndex = index
                this.settingDialog = true
            },
            onShowProjectArticles(index) {
                this.$router.push({name: "Articles", params: {projectId : this.projects[index].projectId, project: this.projects[index]}})
            }
        },
        created: function () {
            if (this.isLogin()) {
                this.isLoading = true
                this.readUserAndProjects(this.$store.getters.email)
            }
        }
    }
</script>


<style>

    @media only screen and (max-width: 768px) {
        .v-content {
            margin: 0;
        }
    }

    .nav-bar {
        margin-bottom: 50px;
    }
</style>
