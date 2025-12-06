import {
  Button,
  Card,
  Checkbox,
  Group,
  Image,
  Loader,
  NumberInput,
  Stack,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { useEffect } from "react";
import { useChannelPointReward } from "../../../hooks/useChannelPointReward";
import type { UserConfig } from "../../../hooks/useUserConfig";
import { type ChannelPointReward, RewardTypes } from "../../../types/Rewards";

const defaultReward: ChannelPointReward = {
  OwnerTwitchID: "",
  Type: RewardTypes.Bttv,
  Title: "BetterTTV Emote",
  Cost: 10000,
  Prompt:
    "Add a BetterTTV emote! In the text field, send a link to the BetterTTV emote. Powered by bot.gempir.com",
  BackgroundColor: "",
  IsMaxPerStreamEnabled: false,
  MaxPerStream: 0,
  IsUserInputRequired: true,
  MaxPerUserPerStream: 0,
  IsMaxPerUserPerStreamEnabled: false,
  IsGlobalCooldownEnabled: false,
  GlobalCooldownSeconds: 0,
  ShouldRedemptionsSkipRequestQueue: false,
  ApproveOnly: false,
  Enabled: false,
  AdditionalOptionsParsed: { Slots: 1 },
};

export function BttvForm({ userConfig }: { userConfig: UserConfig }) {
  const [reward, setReward, deleteReward, errorMessage, loading] =
    useChannelPointReward(
      userConfig?.Protected.CurrentUserID,
      RewardTypes.Bttv,
      defaultReward,
    );

  const form = useForm({
    initialValues: {
      title: reward.Title,
      prompt: reward.Prompt,
      cost: reward.Cost,
      slots: reward.AdditionalOptionsParsed?.Slots || 1,
      backgroundColor: reward.BackgroundColor || "",
      maxPerStream: reward.MaxPerStream,
      maxPerUserPerStream: reward.MaxPerUserPerStream,
      globalCooldownMinutes: reward.GlobalCooldownSeconds / 60,
      approveOnly: reward.ApproveOnly,
      enabled: reward.Enabled,
    },
    validate: {
      title: (value) => (!value ? "Title is required" : null),
      cost: (value) => (value < 1 ? "Cost must be at least 1" : null),
      slots: (value) => (value < 1 ? "Slots must be at least 1" : null),
    },
  });

  useEffect(() => {
    form.setValues({
      title: reward.Title,
      prompt: reward.Prompt,
      cost: reward.Cost,
      slots: reward.AdditionalOptionsParsed?.Slots || 1,
      backgroundColor: reward.BackgroundColor || "",
      maxPerStream: reward.MaxPerStream,
      maxPerUserPerStream: reward.MaxPerUserPerStream,
      globalCooldownMinutes: reward.GlobalCooldownSeconds / 60,
      approveOnly: reward.ApproveOnly,
      enabled: reward.Enabled,
    });
  }, [reward, form.setValues]);

  const handleSubmit = form.onSubmit((values) => {
    const rewardData: ChannelPointReward = {
      OwnerTwitchID: userConfig?.Protected.CurrentUserID,
      Type: RewardTypes.Bttv,
      ApproveOnly: values.approveOnly,
      Title: values.title,
      Prompt: values.prompt,
      Cost: values.cost,
      BackgroundColor: values.backgroundColor,
      IsMaxPerStreamEnabled: Boolean(values.maxPerStream),
      MaxPerStream: values.maxPerStream,
      IsUserInputRequired: true,
      MaxPerUserPerStream: values.maxPerUserPerStream,
      IsMaxPerUserPerStreamEnabled: Boolean(values.maxPerUserPerStream),
      IsGlobalCooldownEnabled: Boolean(values.globalCooldownMinutes),
      GlobalCooldownSeconds: values.globalCooldownMinutes * 60,
      ShouldRedemptionsSkipRequestQueue: false,
      Enabled: values.enabled,
      AdditionalOptionsParsed: {
        Slots: values.slots,
      },
    };

    setReward(rewardData);

    notifications.show({
      title: "Reward Updated",
      message: "BTTV reward settings have been saved",
      color: "green",
    });
  });

  const handleDelete = () => {
    if (confirm("Are you sure you want to delete this reward?")) {
      deleteReward();
      notifications.show({
        title: "Reward Deleted",
        message: "BTTV reward has been removed",
        color: "green",
      });
    }
  };

  if (loading) {
    return (
      <Card shadow="sm" padding="xl" radius="md" withBorder>
        <Group justify="center" p="xl">
          <Loader size="lg" />
        </Group>
      </Card>
    );
  }

  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder>
      <form onSubmit={handleSubmit}>
        <Stack gap="md">
          <Group align="flex-start" wrap="nowrap">
            <Image
              src="https://cdn.betterttv.net/emote/55028cd2135896936880fdd7/3x.webp"
              alt="BTTV"
              w={60}
              h={60}
            />
            <div style={{ flex: 1 }}>
              <Title order={3} size="h4">
                BetterTTV Emotes
              </Title>
              <Text size="sm" c="dimmed">
                Allow viewers to add BTTV emotes via channel points
              </Text>
            </div>
          </Group>

          <TextInput
            label="Reward Title"
            placeholder="BetterTTV Emote"
            required
            {...form.getInputProps("title")}
          />

          <TextInput
            label="Description"
            placeholder="Add a BetterTTV emote..."
            required
            {...form.getInputProps("prompt")}
          />

          <Group grow>
            <NumberInput
              label="Cost (Channel Points)"
              placeholder="10000"
              min={1}
              required
              {...form.getInputProps("cost")}
            />

            <NumberInput
              label="Emote Slots"
              placeholder="1"
              min={1}
              max={100}
              description="How many emotes can be active"
              {...form.getInputProps("slots")}
            />
          </Group>

          <TextInput
            label="Background Color"
            placeholder="#9147FF"
            description="Hex color code for the reward"
            {...form.getInputProps("backgroundColor")}
          />

          <Group grow>
            <NumberInput
              label="Max Per Stream"
              placeholder="0"
              min={0}
              description="0 = unlimited"
              {...form.getInputProps("maxPerStream")}
            />

            <NumberInput
              label="Max Per User Per Stream"
              placeholder="0"
              min={0}
              description="0 = unlimited"
              {...form.getInputProps("maxPerUserPerStream")}
            />
          </Group>

          <NumberInput
            label="Global Cooldown (Minutes)"
            placeholder="0"
            min={0}
            description="0 = no cooldown"
            {...form.getInputProps("globalCooldownMinutes")}
          />

          <Stack gap="xs">
            <Checkbox
              label="Require manual approval"
              description="Redemptions must be manually approved before processing"
              {...form.getInputProps("approveOnly", { type: "checkbox" })}
            />

            <Checkbox
              label="Enabled"
              description="Reward is active and can be redeemed"
              {...form.getInputProps("enabled", { type: "checkbox" })}
            />
          </Stack>

          {errorMessage && (
            <Text c="red" size="sm">
              {errorMessage}
            </Text>
          )}

          <Group justify="space-between" mt="md">
            <Button variant="subtle" color="red" onClick={handleDelete}>
              Delete Reward
            </Button>

            <Button type="submit" color="cyan">
              Save Changes
            </Button>
          </Group>
        </Stack>
      </form>
    </Card>
  );
}
