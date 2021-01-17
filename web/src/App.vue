<template>
  <v-app>
    <v-main>
      <Deposit />
    </v-main>
  </v-app>
</template>

<script>
import Deposit from "./components/Deposit.vue";

export default {
  name: "App",
  components: {
    Deposit
  }
};
</script>
