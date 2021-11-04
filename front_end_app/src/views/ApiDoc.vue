<template>
  <v-app>
    <v-container fluid>
      <v-layout row wrap class="ma-1">
        <v-flex xs12 class="mb-4">
          <div id="button-wrap">
            <v-btn tile dark color="info" @click="expandCollapse">{{expandOrCollapse}}</v-btn>
          </div>
        </v-flex>
        <br><br>
        <v-flex xs12>
          <v-expansion-panels focusable multiple id="panel" v-model="panel">
            <v-expansion-panel v-for="(api, i) in apis" :key="i">
              <v-expansion-panel-header>{{api.title}}</v-expansion-panel-header>
              <v-expansion-panel-content>
                <api-card :apiContent="api"></api-card>
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-flex>
      </v-layout>
    </v-container>
  </v-app>
</template>

<script>
import ApiCard from '../components/ApiCard.vue'

export default {
  name: "ApiDoc",
  components: {ApiCard},
  data() {
    return {
      apis: [
        {
          title: "Adding a subscriber",
          method: "POST",
          path: "/api/v1/subscriber",
          data: [
            {"field": "projectId", "desc": "The unique identifier of the project.", "require": true},
            {"field": "email", "desc": "The pending subscriber's email address.", "require": true},
            {"field": "firstName", "desc": "The pending subscriber's first name.", "require": false},
            {"field": "lastName", "desc": "The pending subscriber's last name.", "require": false}
          ],
          params: []
        },
        {
          title: "Deleting a subscriber",
          method: "DELETE",
          path: "/api/v1/subscriber/:projectId/:email",
          data: [],
          params: [
            {"field": "projectId", "desc": "The unique identifier of the project.", "require": true},
            {"field": "email", "desc": "The subscriber's email address.", "require": true},
          ]
        }
      ],
      panel: [],
      isExpandAll: false
    }
  },
  computed: {
    expandOrCollapse() {
      if (!this.isExpandAll) {
        return "Expand All"
      } else {
        return "Collapse"
      }
    }
  },
  methods: {
    expandALl() {
      this.panel = [...Array(this.apis.length).keys()].map((k, i) => i)
    },
    collapseAll() {
      this.panel = []
    },
    expandCollapse() {
      if (!this.isExpandAll) {
        this.expandALl()
        this.isExpandAll = true
      } else {
        this.collapseAll()
        this.isExpandAll = false
      }
    }
  },
  created() {
  }
}
</script>

<style scoped>
#button-wrap {
  display: inline;
}
#panel {
}
</style>