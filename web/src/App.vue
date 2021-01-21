<template>
  <v-app>
    <v-main>
      <Deposit />
    </v-main>
  </v-app>
</template>

<script>
import Deposit from "./components/Deposit.vue";
import web3 from "./js/web3.js";

export default {
  name: "App",
  components: {
    Deposit
  },
  beforeMount() {
    console.log("Connecting to web3...");
    this.$root.web3 = web3.connectToWeb3();
  }
};
</script>
