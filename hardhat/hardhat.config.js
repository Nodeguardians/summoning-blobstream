require("@nomiclabs/hardhat-ethers");
require("@nomicfoundation/hardhat-chai-matchers");
require("./tasks/proveComet");

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: "0.8.19",
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {},
    sepolia: {
      url: "https://rpc-sepolia-eth.nodeguardians.io",
      accounts: [
        "2897cb159fd9bc3e96fca967549b47268ba1e8d83c59b1419bf2d7ad141ecb40",
      ],
    },
  },
  paths: {
    sources: "./contracts_",
  },
};
