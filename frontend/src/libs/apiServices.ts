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
    console.error("Error in getBlizzardCharacterProfile:", error);
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
