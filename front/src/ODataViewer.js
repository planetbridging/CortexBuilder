import React from "react";
import axios from "axios";
import { Box, Text, Button } from "@chakra-ui/react";
import { FullFileBrowser, setChonkyDefaults } from "chonky";
import { ChonkyIconFA } from "chonky-icon-fontawesome";
import { defineFileAction, ChonkyIconName } from "chonky";

setChonkyDefaults({ iconComponent: ChonkyIconFA });


class ODataViewer extends React.Component {
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
      .get(`http://${this.props.currentHost}:4123/data`)
      .then((response) => {
        const databases = response.data.files.map((db) => ({
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

  listTables = (dbName) => {
    axios
      .get(`http://${this.props.currentHost}:4123/data?dbname=${dbName}`)
      .then((response) => {
        const tables = response.data.files.map((table) => ({
          id: table.name,
          name: table.name,
          isDir: false,
        }));
        this.setState((prevState) => ({
          files: tables,
          currentPath: dbName,
          folderChain: [...prevState.folderChain, { id: dbName, name: dbName }],
        }));
      })
      .catch((error) => console.error("Error fetching tables:", error));
  };

  handleFileAction = (data) => {
    const { action, state } = data;
    const selectedFile = state.selectedFiles[0];

    if (action.id === "open_folder" && selectedFile && selectedFile.isDir) {
      this.listTables(selectedFile.name);
    } else if (
      action.id === "open_folder" &&
      selectedFile &&
      !selectedFile.isDir
    ) {
      //console.log("Selected tbl from db:", selectedFile.name);
    }

    if (action.id === "mount_training_data" && selectedFile) {
      if (!selectedFile.isDir) {
        this.props.handleTrainingDataUpdate(selectedFile.name);
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

    const mountTrainingDataAction = defineFileAction({
      id: "mount_training_data",
      button: {
        name: "Mount Training Data", // Label for the button
        toolbar: false, // Show only in context menu
        contextMenu: true,
        icon: ChonkyIconName.users, // Example icon
      },
      // Disable hotkey for this action
      hotkeys: [],
      // Check if the file is not a directory
      shouldShow: (files) => files.length === 1 && !files[0].isDir,
    });

    return [openFolderAction, mountTrainingDataAction];
  };

  navigateUp = () => {
    this.listDatabases(); // Reset to root view, showing all databases
  };

  render() {
    return (
      <Box h="90%">
        <Text fontSize="3xl">Data Viewer</Text>
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

export default ODataViewer;
