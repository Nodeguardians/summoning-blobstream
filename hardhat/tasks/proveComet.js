task(
  "prove-comet",
  "Prove the comet blob inclusion in the Observatory contract"
)
  .addParam("observatory", "The address of the Observatory contract")
  .addParam("path", "The path to the proof file")
  .setAction(async (taskArgs) => {
    const { observatory, path } = taskArgs;

    let proof = require(path);
    const observatoryContract = await ethers.getContractAt(
      "IObservatory",
      observatory
    );

    const receipt = await observatoryContract.proveComet(proof);
    await receipt.wait();

    const proven = await observatoryContract.isProven();

    if (proven) {
      console.log("Comet is proven!");
    } else {
      console.log("Comet is not proven yet");
    }
  });
