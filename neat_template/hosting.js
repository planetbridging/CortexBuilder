require("dotenv").config();
const express = require("express");
const bodyParser = require("body-parser");
const { MongoClient } = require("mongodb");

function startHosting() {
  const app = express();
  const port = 1789;
  const uri = process.env.MONGODBLINK;
  const client = new MongoClient(uri, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  });

  app.use(bodyParser.json());

  app.post("/saveModel", async (req, res) => {
    try {
      await client.connect();
      const db = client.db(req.body.dbName);
      const collection = db.collection(req.body.collectionName);
      const result = await collection.insertOne(req.body.model);

      res
        .status(201)
        .send({ message: "Model saved successfully", id: result.insertedId });
    } catch (error) {
      console.error("Error saving model:", error);
      res.status(500).send({ message: "Failed to save the model" });
    } finally {
      await client.close();
    }
  });

  app.get("/getModel", async (req, res) => {
    const dbName = req.query.dbName;
    const collectionName = req.query.collectionName;
    const modelId = req.query.modelId; // Assuming you are using _id as a unique identifier

    try {
      await client.connect();
      const db = client.db(dbName);
      const collection = db.collection(collectionName);
      const model = await collection.findOne({
        _id: new MongoClient.ObjectID(modelId),
      });

      if (model) {
        res.status(200).send(model);
      } else {
        res.status(404).send({ message: "Model not found" });
      }
    } catch (error) {
      console.error("Error retrieving model:", error);
      res.status(500).send({ message: "Failed to retrieve the model" });
    } finally {
      await client.close();
    }
  });

  app.get("/listDatabases", async (req, res) => {
    try {
      await client.connect();
      const databasesList = await client.db().admin().listDatabases();
      res.status(200).json(databasesList.databases);
    } catch (error) {
      console.error("Error retrieving databases:", error);
      res.status(500).json({ message: "Failed to retrieve databases" });
    } finally {
      await client.close();
    }
  });

  app.get("/listCollections", async (req, res) => {
    const dbName = req.query.dbName;
    if (!dbName) {
      return res.status(400).json({ message: "Database name is required" });
    }

    try {
      await client.connect();
      const collections = await client.db(dbName).listCollections().toArray();
      res.status(200).json(collections);
    } catch (error) {
      console.error("Error retrieving collections:", error);
      res.status(500).json({ message: "Failed to retrieve collections" });
    } finally {
      await client.close();
    }
  });

  app.get("/listModels", async (req, res) => {
    const dbName = req.query.dbName;
    const collectionName = req.query.collectionName;
    if (!dbName || !collectionName) {
      return res
        .status(400)
        .json({ message: "Database and collection names are required" });
    }

    try {
      await client.connect();
      const collection = client.db(dbName).collection(collectionName);
      const models = await collection
        .find({}, { projection: { _id: 1, name: 1 } })
        .toArray();
      res.status(200).json(models);
    } catch (error) {
      console.error("Error retrieving models:", error);
      res.status(500).json({ message: "Failed to retrieve models" });
    } finally {
      await client.close();
    }
  });

  app.get("/", (req, res) => {
    res.send("Hello, World!");
  });

  app.listen(port, () => {
    console.log(`Server running on port ${port}`);
  });
}

module.exports = {
  startHosting,
};
