<template>
    <StackLayout class="bubble" backgroundColor="#C2CAE8">
        <Label class="user" :text="username"/>
        <Label class="text" :text="message"/>
    </StackLayout>
</template>

<script>
    export default {
        props: {
            username: String,
            message: String
        }
    };
</script>

<style scoped lang="scss">
    // Start custom common variables
    @import '../app-variables';
    // End custom common variables

    // Custom styles
    .fa {
        color: $accent-dark;
    }

    .info {
        font-size: 20;
    }
    .bubble {
        font-size: 20px;
        color: black;
        background-color: gray;
        margin-top: 25px;
        max-width: 80%;
        padding: 15px;
        .personal {
            text-align: right;
        }
        .user {
            align-self: center;
            color: $theme-primary;
        }
        .text {
            color: $theme-secondary;
        }
    }
</style>
