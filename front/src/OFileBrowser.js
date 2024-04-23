import React from "react";
import axios from "axios";
import { FullFileBrowser } from "chonky";
import { FileBrowser, ChonkyIconFA } from "chonky-icon-fontawesome";
import { Box, Text } from "@chakra-ui/react";

class OFileBrowser extends React.Component {
  state = {
    files: [],
    folderChain: [],
  };

  componentDidMount() {
    this.fetchFileData();
  }

  fetchFileData = () => {
    axios
      .get("http://localhost:4123/files") // Change URL as needed
      .then((response) => {
        // Assuming the API returns data directly in the required format
        this.setState({
          files: response.data.files,
          folderChain: response.data.folderChain,
        });
      })
      .catch((error) => {
        console.error("Error fetching file data:", error);
      });
  };

  render() {
    const { files, folderChain } = this.state;

    return (
      <Box h="90%">
        <Text fontSize="3xl">Files</Text>
        <FullFileBrowser
          files={files}
          folderChain={folderChain}
          fileActions={[]}
          iconComponent={ChonkyIconFA}
        />
      </Box>
    );
  }
}

export default OFileBrowser;
