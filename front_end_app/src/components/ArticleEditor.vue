<template>
    <v-main>
        <v-card color="blue-grey darken-1">
            <v-container grid-list-sm>
                <v-layout row wrap>
                    <v-flex xs12>
                        <v-row>
                            <v-btn tile @click="close" class="ma-1">Close</v-btn>
                            <v-spacer></v-spacer>
                            <v-btn tile v-if="article.createdAt > 0" @click="save" class="ma-1 primary">Save</v-btn>
                            <v-btn tile v-else @click="create" class="ma-1 primary">Create</v-btn>
                        </v-row>
                    </v-flex>
                </v-layout>
            </v-container>
            <v-container>
                <v-layout>
                    <v-flex xs12>
                        <v-text-field v-model="article.title" label="Title" filled dark hide-details ></v-text-field>
                    </v-flex>
                </v-layout>
            </v-container>
            <v-container>
                <v-layout row wrap>
                    <v-flex xs12>
                        <tiptap-vuetify v-model="article.htmlBody" :extensions="extensions" placeholder="Start composing..." min-height="900"/>
                    </v-flex>
                </v-layout>
            </v-container>
        </v-card>
        <v-row justify="center">
            <v-dialog v-model="confirmDialog" max-width="250">
                <confirm-dialog v-bind:title="`Are you sure?`" :body="`Your changes will not be saved.`" @yes="confirmedClose" @no="confirmDialog=false"></confirm-dialog>
            </v-dialog>
        </v-row>
    </v-main>
</template>

<script>
    // import the component and the necessary extensions
    import { TiptapVuetify, Heading, Bold, Italic, Strike, Underline, Code, Paragraph, BulletList, OrderedList, ListItem, Link, Blockquote, HorizontalRule, History, CodeBlock, HardBreak, Image } from 'tiptap-vuetify'
    import ConfirmDialog from "./ConfirmDialog";
    import FileSelector from './Gallery'
    import ImageForm from './ImageForm'
    export default {
        // specify TiptapVuetify component in "components"
        components: { TiptapVuetify, ConfirmDialog },
        props: [
            'article'
        ],
        data() {
            return {
                confirmDialog: false,
                titleCopy: "",
                htmlBodyCopy: "",
                extensions: [
                    History,
                    Blockquote,
                    Link,
                    Underline,
                    Strike,
                    Italic,
                    ListItem,
                    BulletList,
                    OrderedList,
                    [Heading, {
                        options: {
                            levels: [1, 2, 3]
                        }
                    }],
                    Bold,
                    Code,
                    CodeBlock,
                    HorizontalRule,
                    Paragraph,
                    HardBreak,
                    [Image, {
                        options: {
                            imageSources: [
                                    { component: FileSelector, name: 'Uploads' },
                                    { component: ImageForm, name: 'Add By Link' }
                                ],
                            imageSourcesOverride: true
                        }
                    }],
                ],
            }
        },
        methods: {
            close() {
                if (this.article.htmlBody != this.htmlBodyCopy || this.article.title != this.titleCopy) {
                    this.confirmDialog = true
                } else {
                    this.confirmedClose()
                }
            },
            confirmedClose() {
                this.$emit("closeEditor")
                this.confirmDialog = false
            },
            save() {
                this.$emit("save", this.article)
            },
            create() {
                this.$emit("create", this.article)
            }
        },
        created() {
            this.titleCopy = this.article.title
            this.htmlBodyCopy = this.article.htmlBody
        }
    }
</script>