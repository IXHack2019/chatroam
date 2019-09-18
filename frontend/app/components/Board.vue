<template>
    <StackLayout class="messages" width="100%" height="100%" orientation="vertical" backgroundColor="lightgray">
        <Message v-bind:username="m.username" v-bind:message="m.message" v-bind:key="m.timestamp" v-for="m in messages"/>
    </StackLayout>
</template>

<script>
    import Message from './Message';

    export default {
        computed: {
            messages() {
                return [ 
                    {username: "jayrsawal", message:"wtf", timestamp: 1234567854321},
                    {username: "dumtard", message:"is going on!!!!", timestamp: 1234567854321} 
                ];
            }
        }
        , components: {
            Message
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
