<template>
    <FlexboxLayout class="messages" width="100%" height="10%" orientation="vertical" backgroundColor="gray">
        <TextField class="input input-border" width="80%"></TextField>
        <Button width="20%" horizontalAlignment="right">&gt;</Button>
    </FlexboxLayout>
</template>

<script>
    import Message from './Message';
    export default {
        
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
