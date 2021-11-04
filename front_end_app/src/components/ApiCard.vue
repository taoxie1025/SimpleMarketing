<template>
    <v-card flat class="mx-auto">
      <v-container fluid>
        <v-layout row wrap>
          <v-card-text>
            <v-flex xs12 sm12>
              <v-chip label :class="getLabel(apiContent.method)">
                {{apiContent.method}}
              </v-chip>
              &nbsp;
              <span class="heading">{{this.fullEndpoint}}</span>
            </v-flex>
            <div v-if="apiContent.data.length > 0">
              <v-divider></v-divider>
              <p>Request Body:</p>
              <v-flex xs12 sm12 v-for="(body, i) in apiContent.data" v-bind:key="i">
                <v-chip ripple label color="white" class="heading">{{body.field}}</v-chip>
                &nbsp;
                <span v-if="body.require" class="require-font">(required) </span>
                <span class="subheading">{{body.desc}}</span>
                <span v-if="body.field=='projectId'">&nbsp;Your projectId is {{projectId}}.</span>
              </v-flex>
            </div>
            <div v-if="apiContent.params.length > 0">
              <v-divider></v-divider>
              <p>Url Params:</p>
              <v-flex xs12 sm12 v-for="(body, i) in apiContent.params" v-bind:key="i">
                <v-chip label color="white" class="heading">:{{body.field}}</v-chip>
                &nbsp;
                <span v-if="body.require" class="require-font">(required) </span>
                <span class="subheading">{{body.desc}}</span>
                <span v-if="body.field=='projectId'">&nbsp;Your projectId is {{projectId}}.</span>
              </v-flex>
            </div>
          </v-card-text>
        </v-layout>
      </v-container>
    </v-card>
</template>

<script>


export default {
  name: "ApiCard",
  components: {},
  props: [
    'apiContent'
  ],
  data() {
    return {
      endpoint: "",
      port: "",
      protocol: ""
    }
  },
  computed: {
    fullEndpoint() {
      return this.protocol + "//" + this.endpoint + (this.port != '80' && this.port != '' ? ":" + this.port : "") + this.apiContent.path
    },
    projectId() {
      return "project-" + window.location.href.split("project-")[1].split("/api")[0]
    }
  },
  methods: {
    getLabel(method) {
      switch (method) {
        case 'POST':
          return 'green'
        case 'DELETE':
          return 'red'
        case 'PUT':
          return 'blue'
        default:
          return 'grey'
      }
    }
  },
  created() {
    this.endpoint = window.location.hostname
    this.port = window.location.port
    this.protocol = window.location.protocol
  }
}
</script>

<style scoped>
.require-font{
  font-style: italic;
}
</style>