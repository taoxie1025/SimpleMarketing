<template>
    <div>
        <v-text-field v-model="form.src" label="link"/>
        <v-text-field v-model="form.alt" label="alt"/>
        <v-btn @click="addImage">
            Add Image
        </v-btn>
    </div>
</template>

<script>
    export default {
        name: "ImageForm",
        components: {  },
        data() {
            // Some public domain images from wikimedia.
            return {
                form: {
                    src: "",
                    alt: ""
                }
            }
        },
        methods: {
            addImage () {
                this.$emit('select-file', {
                    src: this.form.src,
                    alt: this.form.alt
                })
            }
        },
        created() {

        }
    }
</script>

<style scoped>
</style>