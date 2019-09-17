<template>
    <Page class="page" actionBarHidden="true">
        <GridLayout height="100%" width="100%">
            <Board height="90%" verticalAlignment="top"/>
            <Input height="10%" verticalAlignment="bottom"/>
        </GridLayout>
    </Page>
</template>

<script>
    import Board from "./Board";
    import Input from './Input';
    require('nativescript-websockets');
    export default {
        components: {
            Board, Input
        }
    };
    var mySocket = new WebSocket("ws://localhost:8080");
    mySocket.on('open', function(socket) {
        console.log("connected!");
    })
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
