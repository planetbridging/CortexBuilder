import React from "react";
import axios from "axios";
import { Box, Text, Button } from "@chakra-ui/react";
import { FullFileBrowser, setChonkyDefaults } from "chonky";
import { ChonkyIconFA } from "chonky-icon-fontawesome";
import { defineFileAction, ChonkyIconName } from "chonky";

setChonkyDefaults({ iconComponent: ChonkyIconFA });

class OFileBrowser extends React.Component {
  state = {
    files: [],
    folderChain: [{ id: "/", name: "Home" }],
    currentPath: "/", // Start at root
  };

  componentDidMount() {
    this.fetchFileData(this.state.currentPath);
  }

  fetchFileData = (path) => {
    axios
      .get(`http://localhost:4123/files/host${path}`) // Adjust the URL as needed
      .then((response) => {
        this.setState({
          files: response.data.files,
          folderChain: response.data.folderChain,
          currentPath: path,
        });
      })
      .catch((error) => {
        console.error("Error fetching file data:", error);
      });
  };

  handleFileAction = (data) => {
    //console.log(data);

    const fileID = data?.action?.id;
    //console.log(fileID);

    if (fileID == "open_folder") {
      var selectedFolderArray = data?.state?.selectedFiles;
      if (selectedFolderArray) {
        if (selectedFolderArray.length > 0) {
          var tmpKeyCheck = Object.keys(selectedFolderArray[0]);
          if (tmpKeyCheck.includes("isDir")) {
            if (selectedFolderArray[0].isDir == true) {
              //console.log(selectedFolderArray[0].name);
              const tmpNewPath = `${this.state.currentPath}${selectedFolderArray[0].name}/`;
              //console.log(tmpNewPath);
              this.fetchFileData(tmpNewPath);
            }
          }
        }
      }
    }
    /*console.log("Received action:", data.action); // Log the action object for debugging

    const { action, payload } = data;
    const { targetFile } = payload;

    if (action.id === "open_folder" && targetFile.isDir) {
      const newPath = `${this.state.currentPath}${targetFile.id}/`;
      this.fetchFileData(newPath);
    } else if (action.id === "show_file_details" && !targetFile.isDir) {
      alert(`Showing details for ${targetFile.name}`);
    }*/
  };

  setupFileActions = () => {
    const OpenFolderAction = defineFileAction({
      id: "open_folder",
      button: {
        name: "Open Folder",
        toolbar: true,
        contextMenu: true,
        icon: ChonkyIconName.folderOpen,
      },
      hotkeys: ["enter"],
    });

    const ShowFileDetailsAction = defineFileAction({
      id: "show_file_details",
      button: {
        name: "Show Details",
        toolbar: true,
        contextMenu: true,
        icon: ChonkyIconName.fileText,
      },
    });

    return [OpenFolderAction, ShowFileDetailsAction];
  };

  navigateUp = () => {
    const path = this.state.currentPath.split("/").filter(Boolean);
    if (path.length > 0) {
      path.pop(); // Remove last segment of the path
      this.fetchFileData(`/${path.join("/")}/`);
    }
  };

  navigateToProjects = () => {
    const projectsPath = "/projects/"; // Adjust based on your actual path structure
    this.fetchFileData(projectsPath);
  };

  render() {
    const { files, folderChain } = this.state;
    if (!files) {
      return <div>No access?</div>;
    }

    return (
      <Box h="90%">
        <Text fontSize="3xl">Files</Text>
        <Button onClick={this.navigateUp} colorScheme="blue" m={2}>
          Go Up
        </Button>
        <Button onClick={this.navigateToProjects} colorScheme="green" m={2}>
          Go to Projects
        </Button>
        <FullFileBrowser
          files={files}
          folderChain={folderChain}
          fileActions={this.setupFileActions()}
          onFileAction={this.handleFileAction}
          iconComponent={ChonkyIconFA}
        />
      </Box>
    );
  }
}

export default OFileBrowser;
