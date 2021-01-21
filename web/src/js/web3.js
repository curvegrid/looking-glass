// Copyright (c) 2021 Curvegrid Inc.

import { ethers } from "ethers";

export default {
  // connectToWeb3 returns the web3 provider object and the associated configuration details
  connectToWeb3() {
    const result = {
      provider: null,
      available: false
    };
    // Use MetaMask's provider
    result.provider = new ethers.providers.Web3Provider(window.ethereum);
    result.available = true;
    return result;
  },
  formatRawTransaction(rawTX) {
    const tx = JSON.parse(JSON.stringify(rawTX));
    tx.gasLimit = tx.gas;
    tx.gasPrice = ethers.BigNumber.from(tx.gasPrice);
    tx.value = ethers.BigNumber.from(tx.value);
    delete tx.gas;
    delete tx.hash;
    delete tx.from;
    return tx;
  },
  // sign a raw transaction received from the back-end API,
  // then send the signed transaction to the blockchain using web3
  async signRawTransaction(web3, rawTX) {
    const signer = web3.getSigner(rawTX.from);
    const tx = this.formatRawTransaction(rawTX);
    const txResponse = await signer.sendTransaction(tx);
    return txResponse.hash;
  }
};
