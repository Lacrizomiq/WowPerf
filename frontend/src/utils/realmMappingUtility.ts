// services/realmService.ts

import { usRealms } from "./RealmsMapping/usRealmsMapping";
import { euRealms } from "./RealmsMapping/euRealmsMapping";
import { asiaRealms } from "./RealmsMapping/asiaRealmsMapping";

export interface Realm {
  id: number;
  name: string;
  slug: string;
  region: string;
}

export type RegionType = "US" | "EU" | "KR" | "TW";

export const isValidRegion = (region: string): region is RegionType => {
  return ["US", "EU", "KR", "TW"].includes(region.toUpperCase());
};

class RealmService {
  private allRealms: Realm[];

  constructor() {
    // Combine all realm data
    this.allRealms = [...usRealms, ...euRealms, ...asiaRealms];
  }

  // Get realm information by ID
  public getRealmById(realmId: number): Realm | undefined {
    return this.allRealms.find((realm) => realm.id === realmId);
  }

  // Get realm information by name
  public getRealmByName(name: string): Realm | undefined {
    return this.allRealms.find((realm) => realm.name === name);
  }

  // Get realm information by slug
  public getRealmBySlug(slug: string): Realm | undefined {
    return this.allRealms.find((realm) => realm.slug === slug);
  }

  // Build character URL using realm ID
  public buildCharacterUrl(
    characterName: string,
    realmId: number,
    region: RegionType
  ): string {
    const realm = this.getRealmById(realmId);

    if (!realm) {
      console.warn(
        `Realm with ID ${realmId} not found, URL might be incorrect`
      );
      return "";
    }

    return this.formatCharacterUrl(characterName, realm.slug, region);
  }

  // Build character URL using realm name and region
  public buildCharacterUrlByName(
    characterName: string,
    realmName: string,
    region: RegionType
  ): string {
    const realm = this.getRealmByName(realmName);

    if (!realm) {
      console.warn(
        `Realm with name ${realmName} not found, URL might be incorrect`
      );
      return "";
    }

    return this.formatCharacterUrl(characterName, realm.slug, region);
  }

  // Format the final URL
  private formatCharacterUrl(
    characterName: string,
    realmSlug: string,
    region: RegionType
  ): string {
    const formattedCharacterName = characterName.toLowerCase();
    const formattedRegion = region.toLowerCase();

    return `/character/${formattedRegion}/${realmSlug}/${formattedCharacterName}`;
  }

  // Get all realms for a specific region
  public getRealmsByRegion(region: RegionType): Realm[] {
    return this.allRealms.filter((realm) => realm.region === region);
  }

  // Validate if a realm exists in a specific region
  public isValidRealm(realmId: number, region: RegionType): boolean {
    const realm = this.getRealmById(realmId);
    return realm ? realm.region === region : false;
  }
}

// Create a singleton instance
export const realmService = new RealmService();

// Export a function to build character URLs
export const buildCharacterUrl = (
  characterName: string,
  realmId: number,
  region: RegionType
): string => {
  return realmService.buildCharacterUrl(characterName, realmId, region);
};
