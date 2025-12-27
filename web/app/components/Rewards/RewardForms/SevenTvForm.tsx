import {
  Box,
  Button,
  Checkbox,
  Group,
  Image,
  Loader,
  NumberInput,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { useEffect } from "react";
import { useChannelPointReward } from "../../../hooks/useChannelPointReward";
import type { UserConfig } from "../../../hooks/useUserConfig";
import { type ChannelPointReward, RewardTypes } from "../../../types/Rewards";

const defaultReward: ChannelPointReward = {
  OwnerTwitchID: "",
  Type: RewardTypes.SevenTv,
  Title: "7TV Emote",
  Cost: 10000,
  Prompt:
    "Add a 7TV emote! In the text field, send a link to the 7TV emote. Powered by bot.gempir.com",
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

export function SevenTvForm({ userConfig }: { userConfig: UserConfig }) {
  const [reward, setReward, deleteReward, errorMessage, loading] =
    useChannelPointReward(
      userConfig?.Protected.CurrentUserID,
      RewardTypes.SevenTv,
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
      title: (value) => (!value ? "required" : null),
      cost: (value) => (value < 1 ? "min: 1" : null),
      slots: (value) => (value < 1 ? "min: 1" : null),
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
      Type: RewardTypes.SevenTv,
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
      title: "saved",
      message: "7tv reward settings updated",
      color: "green",
    });
  });

  const handleDelete = () => {
    if (confirm("delete this reward?")) {
      deleteReward();
      notifications.show({
        title: "deleted",
        message: "7tv reward removed",
        color: "green",
      });
    }
  };

  if (loading) {
    return (
      <Box
        p="lg"
        style={{
          border: "1px solid var(--border-subtle)",
          backgroundColor: "var(--bg-elevated)",
        }}
      >
        <Group justify="center" p="xl">
          <Loader size="sm" />
        </Group>
      </Box>
    );
  }

  return (
    <Box
      p="md"
      style={{
        border: "1px solid var(--border-subtle)",
        backgroundColor: "var(--bg-elevated)",
      }}
    >
      <form onSubmit={handleSubmit}>
        <Stack gap="md">
          {/* Header */}
          <Group align="flex-start" wrap="nowrap" gap="md">
            <Box
              w={48}
              h={48}
              style={{
                border: "1px solid var(--border-subtle)",
                backgroundColor: "var(--bg-surface)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <Image
                src="https://cdn.7tv.app/emote/01JTTR1J4HTK35S08H272VCWV3/2x.avif"
                alt="7TV"
                w={32}
                h={32}
              />
            </Box>
            <div style={{ flex: 1 }}>
              <Text size="sm" fw={600} ff="monospace" c="white">
                7tv_emotes
              </Text>
              <Text size="xs" c="dimmed" ff="monospace">
                allow viewers to add 7tv emotes via channel points
              </Text>
            </div>
            <Group gap="xs">
              <Box
                className={form.values.enabled ? "status-dot status-online" : "status-dot status-offline"}
              />
              <Text size="xs" ff="monospace" c={form.values.enabled ? "terminal" : "dimmed"}>
                {form.values.enabled ? "enabled" : "disabled"}
              </Text>
            </Group>
          </Group>

          {/* Form Fields */}
          <Box
            p="sm"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-surface)",
            }}
          >
            <Stack gap="sm">
              <TextInput
                label="title"
                placeholder="7TV Emote"
                required
                size="xs"
                {...form.getInputProps("title")}
              />

              <TextInput
                label="description"
                placeholder="Add a 7TV emote..."
                required
                size="xs"
                {...form.getInputProps("prompt")}
              />

              <Group grow>
                <NumberInput
                  label="cost"
                  placeholder="10000"
                  min={1}
                  required
                  size="xs"
                  {...form.getInputProps("cost")}
                />

                <NumberInput
                  label="slots"
                  placeholder="1"
                  min={1}
                  max={100}
                  description="active emote limit"
                  size="xs"
                  {...form.getInputProps("slots")}
                />
              </Group>

              <TextInput
                label="background_color"
                placeholder="#9147FF"
                description="hex color code"
                size="xs"
                {...form.getInputProps("backgroundColor")}
              />
            </Stack>
          </Box>

          {/* Limits Section */}
          <Box
            p="sm"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-surface)",
            }}
          >
            <Text size="xs" fw={600} ff="monospace" c="dimmed" mb="sm" tt="uppercase" style={{ letterSpacing: "0.1em" }}>
              limits
            </Text>
            <Stack gap="sm">
              <Group grow>
                <NumberInput
                  label="max_per_stream"
                  placeholder="0"
                  min={0}
                  description="0 = unlimited"
                  size="xs"
                  {...form.getInputProps("maxPerStream")}
                />

                <NumberInput
                  label="max_per_user"
                  placeholder="0"
                  min={0}
                  description="0 = unlimited"
                  size="xs"
                  {...form.getInputProps("maxPerUserPerStream")}
                />
              </Group>

              <NumberInput
                label="cooldown_minutes"
                placeholder="0"
                min={0}
                description="0 = no cooldown"
                size="xs"
                {...form.getInputProps("globalCooldownMinutes")}
              />
            </Stack>
          </Box>

          {/* Options Section */}
          <Box
            p="sm"
            style={{
              border: "1px solid var(--border-subtle)",
              backgroundColor: "var(--bg-surface)",
            }}
          >
            <Text size="xs" fw={600} ff="monospace" c="dimmed" mb="sm" tt="uppercase" style={{ letterSpacing: "0.1em" }}>
              options
            </Text>
            <Stack gap="xs">
              <Checkbox
                label={
                  <Text size="xs" ff="monospace">
                    require_approval
                  </Text>
                }
                description="manual approval before processing"
                size="xs"
                {...form.getInputProps("approveOnly", { type: "checkbox" })}
              />

              <Checkbox
                label={
                  <Text size="xs" ff="monospace" c={form.values.enabled ? "terminal" : undefined}>
                    enabled
                  </Text>
                }
                description="reward is active"
                size="xs"
                {...form.getInputProps("enabled", { type: "checkbox" })}
              />
            </Stack>
          </Box>

          {errorMessage && (
            <Text c="red" size="xs" ff="monospace">
              error: {errorMessage}
            </Text>
          )}

          <Group justify="space-between">
            <Button
              variant="subtle"
              color="red"
              size="xs"
              onClick={handleDelete}
            >
              delete
            </Button>

            <Button type="submit" color="terminal" size="xs">
              save changes
            </Button>
          </Group>
        </Stack>
      </form>
    </Box>
  );
}
