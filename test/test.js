const protobuf = require("protobufjs");
const WebSocket = require("ws");
const axios = require("axios");
// Load the .proto file and create a Root instance
protobuf.load("../internal/types/transaction.proto", async function (err, root) {
    if (err) {
        throw err;
    }

    // Resolve the 'Transaction' message type
    const Transaction = root.lookupType("Transaction");
    
    // Create a new message instance
    const message = Transaction.create({
        owner: "exampleOwner",
        signature: "exampleSignature",
        type:     root.lookupEnum("TransactionType").values.VerificationTransaction, // Use enum value
        verification_data: {
            attestation: "exampleAttestation",
            cid: "exampleCID",
            datasource: "exampleDataSource",
            timestamp: Date.now()
        }
      
    });

    // Verify the message
    const errMsg = Transaction.verify(message);
    if (errMsg) {
        throw new Error(errMsg);
    }
    
    // Encode the message to a Uint8Array (binary format)
    const buffer = Transaction.encode(message).finish();
    const hexString = Buffer.from(buffer).toString('hex');

    // Log the binary payload
    console.log(buffer);

    const url = "http://" + "localhost" + ":26657/broadcast_tx_commit?tx=" + hexString;
          try {
           
            const data = await axios.request(url);
            console.log("transcation response:", data.data);
          } catch (err) {
            console.error(err?.response?.data ?? err);
          }
});