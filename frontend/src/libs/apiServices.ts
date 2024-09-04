import api from "./api";

export const getBlizzardCharacterProfile = async (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  try {
    const { data } = await api.get(
      `/blizzard/characters/${realmSlug}/${characterName}`,
      {
        params: { region, namespace, locale },
      }
    );
    return data;
  } catch (error) {
    console.error("Error in getBlizzardCharacterProfile:", error);
    throw error;
  }
};

// getBlizzardCharacterMythicPlusBestRuns retrieves the best runs for a character in a specific season
export const getBlizzardCharacterMythicPlusBestRuns = async (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string,
  seasonId: string
) => {
  try {
    const { data } = await api.get(
      `/blizzard/characters/${realmSlug}/${characterName}/mythic-keystone-profile/season/${seasonId}`,
      {
        params: { region, namespace, locale },
      }
    );
    return data;
  } catch (error) {
    console.error("Error in getBlizzardCharacterMythicPlusBestRuns:", error);
    if (error instanceof Error && "response" in error) {
      const axiosError = error as any;
      console.error("Response data:", axiosError.response?.data);
      console.error("Response status:", axiosError.response?.status);
      console.error("Response headers:", axiosError.response?.headers);
    }
    throw error;
  }
};

export const getBlizzardCharacterEquipment = async (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  try {
    const { data } = await api.get(
      `/blizzard/characters/${realmSlug}/${characterName}/equipment`,
      {
        params: { region, namespace, locale },
      }
    );
    return data;
  } catch (error) {
    console.error("Error in getBlizzardCharacterEquipment:", error);
    throw error;
  }
};

export const getBlizzardCharacterSpecializations = async (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  const { data } = await api.get(
    `/blizzard/characters/${realmSlug}/${characterName}/specializations`,
    { params: { region, namespace, locale } }
  );
  return data;
};

// Get all the talents tree for a talent tree ID and spec ID
// Example : for a Shaman with the spec Restauration
export const getBlizzardTalentTree = async (
  talentTreeId: number,
  specId: number,
  region: string,
  namespace: string,
  locale: string
) => {
  const { data } = await api.get(
    `/blizzard/data/talent-tree/${talentTreeId}/playable-specialization/${specId}`,
    {
      params: { region, namespace, locale },
    }
  );
  return data;
};

export const getBlizzardMythicDungeonPerSeason = async (seasonSlug: string) => {
  try {
    const { data } = await api.get(
      `/data/mythic-keystone/season/${seasonSlug}/dungeons`,
      {
        params: {},
      }
    );
    return data;
  } catch (error) {
    console.error("Error in getBlizzardMythicDungeonPerSeason:", error);
    throw error;
  }
};
