import api from "./api";

export const getCharacterProfile = async (
  region: string,
  realm: string,
  name: string,
  fields?: string[]
) => {
  const params = new URLSearchParams({
    region,
    realm,
    name: decodeURIComponent(name),
  });
  if (fields && fields.length > 0) {
    params.append("fields", fields.join(","));
  }

  const { data } = await api.get(`/characters`, { params });

  return data;
};

export const getCharacterGear = async (
  region: string,
  realm: string,
  name: string
) => {
  try {
    const data = await getCharacterProfile(region, realm, name, ["gear"]);
    if (!data || !data.gear) {
      throw new Error("No gear data found in the API response");
    }
    return data;
  } catch (error) {
    console.error("Error in getCharacterGear:", error);
    throw error;
  }
};

export const getCharacterTalents = async (
  region: string,
  realm: string,
  name: string
) => {
  try {
    const data = await getCharacterProfile(region, realm, name, [
      "talents:categorized",
    ]);
    if (!data) {
      throw new Error("No data returned from API");
    }
    if (!data.talents?.categorized?.active && !data.talentLoadout) {
      throw new Error("No talent data found");
    }
    console.log("Talent data received:", data);
    return data;
  } catch (error) {
    console.error("Error fetching talent data:", error);
    throw error;
  }
};

export const getCharacterMythicPlusScores = async (
  region: string,
  realm: string,
  name: string
) => {
  const { data } = await api.get(`/characters/mythic-plus-scores`, {
    params: { region, realm, name },
  });
  return data;
};

export const getCharacterRaidProgression = async (
  region: string,
  realm: string,
  name: string
) => {
  const { data } = await api.get(`/characters/raid-progression`, {
    params: { region, realm, name },
  });
  return data;
};

export const getBlizzardCharacterProfile = async (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  const { data } = await api.get(
    `/blizzard/characters/${realmSlug}/${characterName}`,
    {
      params: { region },
    }
  );
  return data;
};

export const getBlizzardCharacterMythicKeystoneProfile = async (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  const { data } = await api.get(
    `/blizzard/characters/${realmSlug}/${characterName}/mythic-keystone-profile`,
    {
      params: { region },
    }
  );
  return data;
};

export const getBlizzardCharacterEquipment = async (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  const { data } = await api.get(
    `/blizzard/characters/${realmSlug}/${characterName}/equipment`,
    {
      params: { region },
    }
  );
  return data;
};

export const getBlizzardCharacterSpecializations = async (
  region: string,
  realmSlug: string,
  characterName: string
) => {
  const { data } = await api.get(
    `/blizzard/characters/${realmSlug}/${characterName}/specializations`,
    {
      params: { region },
    }
  );
  return data;
};
