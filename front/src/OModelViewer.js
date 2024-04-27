import React from "react";
import axios from "axios";
import { Box, Text, Button } from "@chakra-ui/react";
import { FullFileBrowser, setChonkyDefaults } from "chonky";
import { ChonkyIconFA } from "chonky-icon-fontawesome";
import { defineFileAction, ChonkyIconName } from "chonky";

setChonkyDefaults({ iconComponent: ChonkyIconFA });

const hostIP = "localhost";

class OModelViewer extends React.Component {
  state = {
    files: [],
    folderChain: [{ id: "root", name: "Home" }],
    currentPath: "root", // Start at root which shows databases
  };

  componentDidMount() {
    this.listDatabases();
  }

  listDatabases = () => {
    axios
      .get(`http://${hostIP}:1789/listDatabases`)
      .then((response) => {
        const databases = response.data.map((db) => ({
          id: db.name,
          name: db.name,
          isDir: true,
        }));
        this.setState({
          files: databases,
          currentPath: "root",
          folderChain: [{ id: "root", name: "Home" }],
        });
      })
      .catch((error) => console.error("Error fetching databases:", error));
  };

  listCollections = (dbName) => {
    axios
      .get(`http://${hostIP}:1789/listCollections?dbName=${dbName}`)
      .then((response) => {
        const collections = response.data.map((collection) => ({
          id: collection.name,
          name: collection.name,
          isDir: true,
        }));
        this.setState((prevState) => ({
          files: collections,
          currentPath: dbName,
          folderChain: [...prevState.folderChain, { id: dbName, name: dbName }],
        }));
      })
      .catch((error) => console.error("Error fetching collections:", error));
  };

  listModels = (dbName, collectionName) => {
    axios
      .get(
        `http://${hostIP}:1789/listModels?dbName=${dbName}&collectionName=${collectionName}`
      )
      .then((response) => {
        const models = response.data.map((model) => ({
          id: model._id,
          name: model._id,
          isDir: false,
        }));
        this.setState((prevState) => ({
          files: models,
          currentPath: `${dbName}/${collectionName}`,
          folderChain: [
            ...prevState.folderChain,
            { id: collectionName, name: collectionName },
          ],
        }));
      })
      .catch((error) => console.error("Error fetching models:", error));
  };

  getModel = (dbName, collectionName, modelId) => {
    axios
      .get(
        `http://${hostIP}:1789/getModel?dbName=${dbName}&collectionName=${collectionName}&modelId=${modelId}`
      )
      .then((response) => {
        console.log("Model data:", response.data);
      })
      .catch((error) => console.error("Error fetching model:", error));
  };

  handleFileAction = (data) => {
    //console.log("Action data:", data);
    const { action, state } = data;
    const selectedFile = state.selectedFiles[0];
    //console.log("Selected file:", selectedFile);

    if (action.id === "open_folder" && selectedFile && selectedFile.isDir) {
      const { id, name } = selectedFile;
      const { currentPath } = this.state;

      if (currentPath === "root") {
        this.listCollections(name);
      } else {
        const pathParts = currentPath.split("/");
        console.log(pathParts.length);
        if (pathParts.length === 1) {
          this.listModels(pathParts[0], name);
        }

        /*else if (pathParts.length === 2) {
          console.log("hello");
          //this.getModel(pathParts[0], pathParts[1], id);
        }*/
      }
    }
  };

  setupFileActions = () => {
    const openFolderAction = defineFileAction({
      id: "open_folder",
      button: {
        name: "Open",
        toolbar: true,
        contextMenu: true,
        icon: ChonkyIconName.folderOpen,
      },
      hotkeys: ["enter"],
    });

    return [openFolderAction];
  };

  navigateUp = () => {
    const pathParts = this.state.currentPath.split("/").filter(Boolean);
    if (pathParts.length === 1) {
      this.listDatabases();
    } else if (pathParts.length === 2) {
      this.listCollections(pathParts[0]);
    }
  };

  render() {
    return (
      <Box h="90%">
        <Text fontSize="3xl">Model Viewer</Text>
        <Button onClick={this.navigateUp} colorScheme="blue" m={2}>
          Go Up
        </Button>
        <FullFileBrowser
          files={this.state.files}
          folderChain={this.state.folderChain}
          fileActions={this.setupFileActions()}
          onFileAction={this.handleFileAction}
          iconComponent={ChonkyIconFA}
        />
      </Box>
    );
  }
}

export default OModelViewer;
