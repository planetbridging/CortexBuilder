import React, { Component } from "react";
import {
  Box,
  Flex,
  Wrap,
  WrapItem,
  Icon,
  Text,
  Button,
} from "@chakra-ui/react";
import {
  FiFolder,
  FiSettings,
  FiBookOpen,
  FiLogIn,
  FiLogOut,
  FiEye,
} from "react-icons/fi";

class OHome extends Component {
  state = {
    isLoggedIn: false,
    content: "Please log in.",
  };

  handleLoginLogout = () => {
    this.setState((prevState) => ({
      isLoggedIn: !prevState.isLoggedIn,
      content: prevState.isLoggedIn ? "Please log in." : "Welcome back!",
    }));
  };

  render() {
    return (
      <Box display="flex" flexDirection="column" h="100vh" w="100vw">
        {/* Top Menu */}
        <Flex
          justifyContent="space-between"
          alignItems="center"
          bg="gray.100"
          p={4}
          wrap="wrap"
        >
          {/* Logo on the left */}
          <Text fontSize="xl" fontWeight="bold">
            Cortex Builder
          </Text>

          {/* Center Menu */}
          <Wrap spacing={8} wrap="wrap" justify="center">
            <WrapItem>
              <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                <Text mt={2}>
                  <Icon as={FiFolder} /> Files
                </Text>
              </Box>
            </WrapItem>
            <WrapItem>
              <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                <Text mt={2}>
                  <Icon as={FiEye} /> Model Viewer
                </Text>
              </Box>
            </WrapItem>
            <WrapItem>
              <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                <Text mt={2}>
                  <Icon as={FiBookOpen} /> Training
                </Text>
              </Box>
            </WrapItem>
            <WrapItem>
              <Box p={2} bg="white" boxShadow="sm" borderRadius="md">
                <Text mt={2}>
                  <Icon as={FiSettings} /> Settings
                </Text>
              </Box>
            </WrapItem>
          </Wrap>

          {/* Login/Logout on the right */}
          <Button
            onClick={this.handleLoginLogout}
            leftIcon={
              this.state.isLoggedIn ? (
                <Icon as={FiLogOut} />
              ) : (
                <Icon as={FiLogIn} />
              )
            }
          >
            {this.state.isLoggedIn ? "Logout" : "Login"}
          </Button>
        </Flex>

        {/* Main content */}
        <Box flex="1" p={4} bg="gray.50">
          {/* Dynamic content based on login status */}
          {this.state.content}
        </Box>
      </Box>
    );
  }
}

export default OHome;
