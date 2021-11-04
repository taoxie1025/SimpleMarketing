<template>
    <v-card>
        <v-card-title v-if="title" class="headline">{{title}}</v-card-title>
        <v-card-text v-if="body">{{body}}</v-card-text>
        <v-card-actions>
            <v-spacer></v-spacer>
            <div v-if="!simpleConfirm">
                <v-btn color="red darken-1" text @click="yes()">Yes</v-btn>
                <v-btn color="green darken-1" text @click="no()">No</v-btn>
            </div>
            <div v-else>
                <v-btn color="red darken-1" text @click="no()">Close</v-btn>
            </div>
        </v-card-actions>
    </v-card>
    </template>
<script>

    export default {
        data() {
            return {
            }
        },
        props:[
            'title',
            'body',
            'simpleConfirm'
        ],
        methods: {
            yes() {
                this.$emit("yes")
            },
            no() {
                this.$emit("no")
            }
        }
    }
</script>