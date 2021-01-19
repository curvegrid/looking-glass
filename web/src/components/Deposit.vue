<template>
  <div>
    <v-form ref="form">
      <v-container>
        <v-text-field v-model="amount" label="Amount" />
        <v-text-field v-model="recipient" label="Recipient" />
        <v-autocomplete
          v-model="originToken"
          :items="tokens"
          :loading="isLoading"
        />
        <v-autocomplete
          v-model="destinationToken"
          :items="tokens"
          :loading="isLoading"
        />
        <v-btn @click="deposit" :disabled="disableFunctionSubmit">
          Deposit
        </v-btn>
      </v-container>
    </v-form>
  </div>
</template>

<script>
// Copyright (c) 2021 Curvegrid Inc.

import axios from "axios";
import web3 from "../js/web3.js";

export default {
  data() {
    return {
      // form data
      amount: 0,
      recipient: "",
      originToken: {},
      destinationToken: {},

      // program state
      isLoading: false,

      // parsed data
      tokens: []
    };
  },
  computed: {
    disableFunctionSubmit() {
      return !this.$root.web3.available;
    }
  },
  async created() {
    await this.getTokens();
  },
  methods: {
    async getTokens() {
      try {
        this.isLoading = true;
        const response = await axios.get("/api/resources");
        const resourceMapping = response.data;
        this.parseTokens(resourceMapping);
      } catch (err) {
        console.error(err);
      }
      this.isLoading = false;
    },
    // parse a mapping from resource id to a list of
    // corresponding resources
    // into a list of tokens
    // Note: resource's structure: ({ chainID, tokenAddress, tokenHandlerAddress })
    async parseTokens(resourceMapping) {
      this.tokens = [];
      Object.values(resourceMapping).forEach(resources => {
        resources.forEach(resource => {
          this.tokens.push({
            value: {
              tokenAddress: resource.tokenAddress,
              chainID: resource.chainID
            },
            text: `address: ${resource.tokenAddress}, chainID: ${resource.chainID}`
          });
        });
      });
    },
    async deposit() {
      const depositData = {
        amount: this.amount,
        recipient: this.recipient,
        originChainID: this.originToken.chainID,
        originTokenAddress: this.originToken.tokenAddress,
        destinationChainID: this.destinationToken.chainID,
        destinationTokenAddress: this.destinationToken.tokenAddress
      };
      try {
        const response = await axios.post("/api/deposit", depositData);
        const tx = response.data.tx;
        if (this.$root.web3.available)
          await web3.signRawTransaction(this.$root.web3.provider, tx);
        else
          throw new Error(
            `web3 is not available, cannot sign transaction: ${tx}`
          );
      } catch (err) {
        console.error(err);
      }
    }
  }
};
</script>
