export type Badge = {
  id: string;
  name: string;
  description: string;
  imageUrl: string;
  conditionType: string;
  conditionId: string;
};

export type UserBadge = {
  id: string;
  userId: string;
  badge: Badge;
  earnedAt: string;
};

export type UserBadgesResponse = {
  badges: UserBadge[];
};

export type CompleteLessonResult = {
  progress: {
    id: string;
    userId: string;
    lessonId: string;
    completedAt: string;
  };
  newBadges: UserBadge[];
};
