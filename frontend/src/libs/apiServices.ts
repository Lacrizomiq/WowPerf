import api from "./api";
import axios from "axios";
import { StaticRaid } from "@/types/raids";

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
  if (seasonId === "13") {
    // TWW Season 1
    console.log("Using static data for TWW Season 1");
    return null; // Return null for TWW Season 1
  }

  try {
    const { data } = await api.get(
      `/blizzard/characters/${realmSlug}/${characterName}/mythic-keystone-profile/season/${seasonId}`,
      {
        params: { region, namespace, locale },
      }
    );
    return data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      if (error.response?.status === 404) {
        console.log(`No Mythic+ data available for season ${seasonId}`);
        return null;
      }
      console.error(
        "Error in getBlizzardCharacterMythicPlusBestRuns:",
        error.message
      );
    }
    console.error(
      "Unexpected error in getBlizzardCharacterMythicPlusBestRuns:",
      error
    );
    return null;
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

export const getBlizzardRaidsByExpansion = async (
  expansion: string
): Promise<StaticRaid[]> => {
  try {
    const url = `/blizzard/data/raids/${expansion}`;
    console.log(`Fetching raids for expansion: ${expansion} from URL: ${url}`);
    const { data } = await api.get<StaticRaid[]>(url);
    return data;
  } catch (error) {
    console.error("Error in getBlizzardRaidsByExpansion:", error);
    throw error;
  }
};

export const getBlizzardCharacterEncounterRaid = async (
  region: string,
  realmSlug: string,
  characterName: string,
  namespace: string,
  locale: string
) => {
  try {
    const url = `/blizzard/characters/${realmSlug}/${characterName}/encounters/raids`;
    console.log(`Fetching raid encounters from: ${url}`);
    console.log(
      `Params: region=${region}, namespace=${namespace}, locale=${locale}`
    );
    const { data } = await api.get(url, {
      params: { region, namespace, locale },
    });
    console.log("Raid encounter data received:", data);
    return data;
  } catch (error: any) {
    if (error.response && error.response.status === 404) {
      console.warn("Raid encounter data not found for this character");
      return null;
    }
    console.error("Error in getBlizzardCharacterEncounterRaid:", error);
    throw error;
  }
};
