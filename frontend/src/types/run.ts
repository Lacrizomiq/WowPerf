export interface Run {
  clear_time_ms: number;
  completed_at: string;
  deleted_at: null | string;
  dungeon: {
    expansion_id: number;
    group_finder_activity_ids: number[];
    icon_url: string;
    id: number;
    keystone_timer_ms: number;
    map_challenge_mode_id: number;
    name: string;
    num_bosses: number;
    patch: string;
    short_name: string;
    slug: string;
    type: string;
    wowInstanceId: number;
  };
  faction: string;
  keystone_platoon_id: null | number;
  keystone_run_id: number;
  keystone_team_id: number;
  keystone_time_ms: number;
  logged_run_id: null | number;
  mythic_level: number;
  num_chests: number;
  num_modifiers_active: number;
  platoon: null | any;
  roster: Array<{
    character: {
      class: {
        id: number;
        name: string;
        slug: string;
      };
      faction: string;
      id: number;
      level: number;
      name: string;
      path: string;
      persona_id: number;
      race: {
        faction: string;
        id: number;
        name: string;
        slug: string;
      };
      realm: {
        altName: null | string;
        altSlug: string;
        connectedRealmId: number;
        id: number;
        isConnected: boolean;
        locale: string;
        name: string;
        realmType: string;
        slug: string;
        wowConnectedRealmId: number;
        wowRealmId: number;
      };
      recruitmentProfiles: any[];
      region: {
        name: string;
        short_name: string;
        slug: string;
      };
      spec: {
        id: number;
        name: string;
        slug: string;
      };
      stream: {
        community_ids: any[];
        game_id: string;
        id: string;
        language: string;
        name: string;
        started_at: string;
        thumbnail_url: string;
        title: string;
        type: string;
        user_id: string;
        viewer_count: number;
      } | null;
    };
    isTransfer: boolean;
    oldCharacter: null | any;
    role: string;
  }>;
  season: string;
  status: string;
  time_remaining_ms: number;
  videos: any[];
  weekly_modifiers: Array<{
    description: string;
    icon: string;
    id: number;
    name: string;
    slug: string;
  }>;
}

export interface MythicPlusRun {
  rank: number;
  run: Run;
  score: number;
}

export interface MythicPlusData {
  leaderboard_url: string;
  params: {
    dungeon: string;
    page: number;
    region: string;
    season: string;
  };
  rankings: MythicPlusRun[];
}
