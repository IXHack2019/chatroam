<template>
    <FlexboxLayout backgroundColor="green">
        <Label :text="username" width="100%" height="50" />
        <Label :text="message" width="100%" height="50" />
    </FlexboxLayout>
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
</style>
