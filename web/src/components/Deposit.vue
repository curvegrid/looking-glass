<template>
  <div>
    <v-form ref="form">
      <v-container>
        <v-text-field label="Amount" />
      </v-container>
    </v-form>
  </div>
</template>

<script>
// Copyright (c) 2021 Curvegrid Inc.

import axios from "axios";

export default {
  data() {
    return {
      amount: 0,
      tokens: []
    };
  },
  async created() {
    await this.getTokens();
  },
  methods: {
    async getTokens() {
      try {
        const response = await axios.get("/api/resources");
        this.resources = response.data;
        console.log(JSON.stringify(this.resources));
      } catch (err) {
        console.error(err);
      }
    }
  }
};
</script>
