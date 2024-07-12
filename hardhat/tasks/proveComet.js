task(
  "prove-comet",
  "Prove the comet blob inclusion in the Observatory contract"
)
  .addParam("observatory", "The address of the Observatory contract")
  .addParam("path", "The path to the proof file")
  .setAction(async (taskArgs) => {
    const { observatory, path } = taskArgs;

    const proof = require(path);
    const observatoryContract = await ethers.getContractAt(
      "IObservatory",
      observatory
    );

    const tx = await observatoryContract.proveComet(proof, "0x00");
    await tx.wait();

    const proven = await observatory.isProven();

    if (proven) {
      console.log("Comet is proven!");
    } else {
      console.log("Comet is not proven yet");
    }
  });
