import React from "react";
import axios from "axios";
import { Box, Text, Button } from "@chakra-ui/react";
import { FullFileBrowser, setChonkyDefaults } from "chonky";
import { ChonkyIconFA } from "chonky-icon-fontawesome";
import { defineFileAction, ChonkyIconName } from "chonky";

import { OFFNN } from "./OFFNN";

setChonkyDefaults({ iconComponent: ChonkyIconFA });



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
      .get(`http://${this.props.currentHost}:1789/listDatabases`)
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
      .get(`http://${this.props.currentHost}:1789/listCollections?dbName=${dbName}`)
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
        `http://${this.props.currentHost}:1789/listModels?dbName=${dbName}&collectionName=${collectionName}`
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
        `http://${this.props.currentHost}:1789/getModel?dbName=${dbName}&collectionName=${collectionName}&modelId=${modelId}`
      )
      .then((response) => {
        console.log("Model data:", response.data);
        console.log("--------running model testing-------");
        const nn = new OFFNN(response.data);

        // Example input values - adjust based on your actual input configuration
        const inputValues = { 1: 1, 2: 0.5, 3: 0.75 };
        const outputs = nn.feedforward(inputValues);
        console.log("Network outputs:", outputs);
      })
      .catch((error) => console.error("Error fetching model:", error));
  };

  handleFileAction = (data) => {
    //console.log("Action data:", data);
    const { action, state } = data;
    const selectedFile = state.selectedFiles[0];
    //console.log("Selected file:", selectedFile);

    if (action.id === "open_folder" && selectedFile) {
      const { id, name } = selectedFile;
      const { currentPath } = this.state;

      if (selectedFile.isDir) {
        if (currentPath === "root") {
          this.listCollections(name);
        } else {
          const pathParts = currentPath.split("/");
          if (pathParts.length === 1) {
            this.listModels(pathParts[0], name);
          }
        }
      } else {
        const pathParts = currentPath.split("/");
        if (pathParts.length === 2) {
          this.getModel(pathParts[0], pathParts[1], id);
        }
      }
    }

    if (action.id === "mount_population" && selectedFile) {
      const { id, name } = selectedFile;
      const { currentPath } = this.state;

      const pathParts = currentPath.split("/");
      const dbName = pathParts[0];
      const collectionName = pathParts[1];


      if(this.props.onDataUpdate){
        this.props.onDataUpdate(dbName,collectionName);
      }

      // Log the required data
      /*console.log(
        `Mount Population clicked for file: ${name} (ID: ${id})`,
        `Database: ${dbName}`,
        `Collection: ${collectionName}`
      );*/
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

    const mountPopulationAction = defineFileAction({
      id: "mount_population",
      button: {
        name: "Mount Population", // Label for the button
        toolbar: false, // Show only in context menu
        contextMenu: true,
        icon: ChonkyIconName.users, // Example icon 
      },
      // Disable hotkey for this action
      hotkeys: [],
      // Check if the file is not a directory
      shouldShow: (files) => files.length === 1 && !files[0].isDir,
    });

    return [openFolderAction,mountPopulationAction];
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
