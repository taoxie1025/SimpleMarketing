<template>
    <v-dialog max-width="600px" v-model="dialog">

        <v-card>
            <v-card-title>
                <h2>Add a New Project</h2>
            </v-card-title>

            <v-card-text>
                <v-form class="px-3" ref="form">

                    <v-alert type="error" :value="alerta">
                        Type some title
                    </v-alert>

                    <v-text-field label="Title" v-model="title" prepend-icon="edit"></v-text-field>

                    <v-spacer></v-spacer>

                    <v-btn text class="success mx-0 mt-3" @click="submit" :loading="loading">Add project</v-btn>
                </v-form>
            </v-card-text>

        </v-card>
    </v-dialog>
</template>

<script>
    import { EventBus } from '@/event_bus.js';
    export default {
        data() {
            return {
                dialog: false,
                title: '',
                content: '',
                due: null,
                userId: '',
                alerta: false,
                dateRules: [
                    v =>(v.length >= 4) || 'Invalid format. Click to pick'
                ],
                loading: false,
            }
        },
        methods: {

            submit() {

            },
            formReset() {
                this.title = '';
                this.content = '';
                this.due = null;
            },
        },
        computed: {

        },
        mounted() {
            EventBus.$on('dialog', (value) => {
                this.dialog = value
            })
        },
        created() {

        }
    }
</script>