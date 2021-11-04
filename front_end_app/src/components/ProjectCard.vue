<template>
  <div>
    <v-container>
      <v-layout>
        <v-flex xs12 sm12>
          <v-hover v-slot:default="{ hover }" open-delay="0">
            <v-card :elevation="hover ? 20 : 4" min-width="330px">
              <v-toolbar extended :src="project.backgroundImageUrl!='' ? project.backgroundImageUrl: defaultBackground[0].icon">
                <template v-slot:extension>
                  <v-btn class="mt-3 ml-3" fab x-large light :disabled="true" id="avatar-button">
                    <avatar v-if="project.avatarUrl" alt="Avatar" :src="project.avatarUrl" :size="69"></avatar>
                    <avatar v-else :username="project.name" backgroundColor="grey" color="white"  :size="69"></avatar>
                  </v-btn>
                </template>
              </v-toolbar>
              <v-card-title>
                <v-row class="ml-0">
                  <v-list-item-content>
                    <v-list-item-title v-html="project.name" class="title-text"></v-list-item-title>
                    <v-list-item-subtitle v-if="project.lastBroadcastTimeMs > 0" class="subtitle-text" >Last broadcast date: {{lastBroadcastTime}}</v-list-item-subtitle>
                    <v-list-item-subtitle v-else class="subtitle-text" >Created on: {{createdTime}} </v-list-item-subtitle>
                  </v-list-item-content>
                  <v-menu right class="ml-0">
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn v-bind="attrs" v-on="on" icon ripple class="ma-2" id="more-btn" >
                        <v-tooltip bottom>
                          <template v-slot:activator="{ on, attrs }">
                            <v-icon size="30px" v-bind="attrs" v-on="on" >mdi-dots-vertical</v-icon>
                          </template>
                          <span>More</span>
                        </v-tooltip>
                      </v-btn>
                    </template>
                    <v-list dense>
                      <v-list-item @click="showSubscribers">
                        <v-list-item-icon>
                          <v-icon color="grey">mdi-account-group-outline</v-icon>
                        </v-list-item-icon>
                        <v-list-item-content>
                          <v-list-item-title>Subscribers</v-list-item-title>
                        </v-list-item-content>
                      </v-list-item>
                      <v-list-item @click="showApiDoc">
                        <v-list-item-icon>
                          <v-icon color="grey">mdi-api</v-icon>
                        </v-list-item-icon>
                        <v-list-item-content>
                          <v-list-item-title>API Doc</v-list-item-title>
                        </v-list-item-content>
                      </v-list-item>
                      <v-list-item @click="showSettings">
                        <v-list-item-icon>
                          <v-icon color="grey">mdi-cog-outline</v-icon>
                        </v-list-item-icon>
                        <v-list-item-content>
                          <v-list-item-title>Settings</v-list-item-title>
                        </v-list-item-content>
                      </v-list-item>
                    </v-list>
                  </v-menu>
                </v-row>
              </v-card-title>
              <v-flex>
                <v-chip class="ml-3" color="transparent" label disabled>
                  <v-icon size="15">mdi-email-send-outline</v-icon>
                  <span class="subheading">Total sent:{{project.totalBroadcastCount}}</span>
                </v-chip>
              </v-flex>
              <v-card-actions>
                <div v-if="project.projectState==3">
                  <v-chip color="transparent" label disabled>
                    <v-icon size="20" style="color: limegreen;">mdi-brightness-1</v-icon>
                    <span class="subheading">Online</span>
                  </v-chip>
                </div>
                <div v-else>
                  <v-chip class="ma-0" color="transparent" label disabled>
                    <v-icon size="20" style="color: red;">mdi-brightness-1</v-icon>
                    <span class="subheading">Offline</span>
                  </v-chip>
                </div>
                <v-spacer></v-spacer>
                <v-btn small tile text color="primary" @click="showArticles">ARTICLES</v-btn>
                <v-spacer></v-spacer>
                <v-btn small tile text outlined color="primary"  @click="goToLandingPage">LANDING PAGE</v-btn>
              </v-card-actions>
            </v-card>
          </v-hover>
        </v-flex>
      </v-layout>
    </v-container>
  </div>
</template>

<script>
    import Avatar from "vue-avatar";

    export default {
      components: {Avatar},
      data() {
            return {
              defaultBackground: [
                { title: 'Dashboard', icon: require('@/assets/images/landing/Bg/pattern.jpg') }
              ],
            }
        },
        props:[
            'project',
            'index'
        ],
        methods: {
            showSubscribers() {
              this.$router.push({path: 'subscribers/' + this.project.projectId});
            },
            showSettings() {
                this.$emit("showProjectSettings", this.index)
            },
            showApiDoc() {
              const routeData = this.$router.resolve({path: 'documents/' + this.project.projectId + '/api'});
              window.open(routeData.href, '_blank');
            },
            showArticles() {
                this.$emit("showProjectArticles", this.index)
            },
            goToLandingPage() {
              const projectIdShort = this.project.projectId.substring('project-'.length, this.project.projectId.length)
              const routeData = this.$router.resolve({path: 'p/' + projectIdShort});
              window.open(routeData.href, '_blank');
            }
        },
        computed: {
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
        mounted() {

        },
        created() {
        }
    }
</script>

<style lang="scss">
    .title-text {
        font-weight: 500;
        font-size: 120%;
        max-width: 300px;
    }
    .subtitle-text {
        font-weight: 300;
        font-size: 50%;
    }
    i.v-icon.v-icon {
        color: red;
    }
    .bottom-card-actions {
        bottom: 0;
        position: absolute;
    }
    .mx-auto {
        margin-bottom: 0px;
    }
    .status-label-actions {
        position:absolute;
        bottom:30px;
    }
    .actions-container {
        margin-top: 10px;
        position:relative;
        bottom: -15px;
        margin-left: 5px;
        margin-right: 5px;
    }
    .title-container {
        margin: 5px;
        display: flex;
        max-width: 360px;
    }
    .bottom-action-container {
        margin-top: 5px;
        margin-left: 5px;
        margin-right: 5px;
    }
    .add-button {
        margin: 0;
        position: absolute;
        top: 50%;
        left: 50%;
        -ms-transform: translate(-50%, -50%);
        transform: translate(-50%, -50%);
    }
    #more-btn {
      outline: none !important;
      box-shadow: none;
    }
    .v-ripple__container {
      display:none !important;
    }
    #avatar-button.v-btn--disabled {
      background-color: white !important;
    }
</style>