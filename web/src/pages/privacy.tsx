import {
  Anchor,
  Container,
  Divider,
  List,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { GetServerSidePropsContext } from "next";
import { initializeStore } from "../service/initializeStore";

export default function Privacy() {
  return (
    <Container size="md">
      <Stack gap="xl" py="xl">
        <div>
          <Title order={1} mb="md">
            Privacy Policy
          </Title>
          <Text c="dimmed">Last updated: 2024</Text>
        </div>

        <Divider />

        <Stack gap="lg">
          <div>
            <Title order={2} size="h3" mb="md">
              Data Collection
            </Title>
            <Text>
              gempbot collects and stores the following information when you use
              our service:
            </Text>
            <List mt="md" spacing="sm">
              <List.Item>Your Twitch user ID and username</List.Item>
              <List.Item>
                Channel point reward configurations you create
              </List.Item>
              <List.Item>
                Emote history and blocked emotes for your channel
              </List.Item>
              <List.Item>User permissions you grant to other users</List.Item>
              <List.Item>Overlay data created through our editor</List.Item>
            </List>
          </div>

          <div>
            <Title order={2} size="h3" mb="md">
              How We Use Your Data
            </Title>
            <Text>
              Your data is used exclusively to provide the bot and overlay
              management services:
            </Text>
            <List mt="md" spacing="sm">
              <List.Item>
                To manage channel point rewards and emote redemptions
              </List.Item>
              <List.Item>
                To store and serve your custom overlay configurations
              </List.Item>
              <List.Item>
                To authenticate your access to the dashboard
              </List.Item>
              <List.Item>To enforce permissions and access controls</List.Item>
            </List>
          </div>

          <div>
            <Title order={2} size="h3" mb="md">
              Third-Party Services
            </Title>
            <Text>
              gempbot integrates with the following third-party services:
            </Text>
            <List mt="md" spacing="sm">
              <List.Item>
                <strong>Twitch:</strong> For authentication and channel point
                reward management
              </List.Item>
              <List.Item>
                <strong>BetterTTV:</strong> For emote management functionality
              </List.Item>
              <List.Item>
                <strong>7TV:</strong> For emote management functionality
              </List.Item>
            </List>
            <Text mt="md" c="dimmed">
              These services may collect their own data according to their
              respective privacy policies.
            </Text>
          </div>

          <div>
            <Title order={2} size="h3" mb="md">
              Data Storage
            </Title>
            <Text>
              All data is stored securely in our database. We implement
              appropriate technical and organizational measures to protect your
              personal information against unauthorized access, alteration,
              disclosure, or destruction.
            </Text>
          </div>

          <div>
            <Title order={2} size="h3" mb="md">
              Data Deletion
            </Title>
            <Text>
              You can request deletion of your data at any time by contacting
              us. When you delete a channel point reward or overlay, the
              associated data is immediately removed from our systems.
            </Text>
          </div>

          <div>
            <Title order={2} size="h3" mb="md">
              Cookies
            </Title>
            <Text>
              gempbot uses cookies to maintain your authentication session.
              These cookies are essential for the service to function and are
              not used for tracking or analytics purposes.
            </Text>
          </div>

          <div>
            <Title order={2} size="h3" mb="md">
              Changes to This Policy
            </Title>
            <Text>
              We may update this privacy policy from time to time. We will
              notify you of any changes by posting the new privacy policy on
              this page and updating the "Last updated" date.
            </Text>
          </div>

          <div>
            <Title order={2} size="h3" mb="md">
              Contact
            </Title>
            <Text>
              If you have any questions about this privacy policy or our data
              practices, please contact us through{" "}
              <Anchor
                href="https://twitch.tv/gempir"
                target="_blank"
                rel="noopener noreferrer"
              >
                Twitch
              </Anchor>
              .
            </Text>
          </div>
        </Stack>
      </Stack>
    </Container>
  );
}

export async function getServerSideProps(ctx: GetServerSidePropsContext) {
  return initializeStore(ctx);
}
