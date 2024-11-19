package data

// GenshinDataFileName represents valid Excel config data filenames
type GenshinDataFileName string

// Character-related files
const (
	CharacterDataFile          GenshinDataFileName = "AvatarExcelConfigData"
	CharacterProfileFile       GenshinDataFileName = "FetterInfoExcelConfigData"
	CharacterCostumeFile       GenshinDataFileName = "AvatarCostumeExcelConfigData"
	CharacterSkillDepotFile    GenshinDataFileName = "AvatarSkillDepotExcelConfigData"
	CharacterSkillFile         GenshinDataFileName = "AvatarSkillExcelConfigData"
	CharacterTalentFile        GenshinDataFileName = "ProudSkillExcelConfigData"
	CharacterConstellationFile GenshinDataFileName = "AvatarTalentExcelConfigData"
	CharacterAscensionFile     GenshinDataFileName = "AvatarPromoteExcelConfigData"
	CharacterStatCurveFile     GenshinDataFileName = "AvatarCurveExcelConfigData"
	CharacterReleaseInfoFile   GenshinDataFileName = "AvatarCodexExcelConfigData"

	ArtifactSetBonusFile   GenshinDataFileName = "EquipAffixExcelConfigData"
	ArtifactDataFile       GenshinDataFileName = "ReliquaryExcelConfigData"
	ArtifactMainStatFile   GenshinDataFileName = "ReliquaryLevelExcelConfigData"
	ArtifactSubStatFile    GenshinDataFileName = "ReliquaryAffixExcelConfigData"
	ArtifactSetDataFile    GenshinDataFileName = "ReliquarySetExcelConfigData"
	ArtifactRarityDataFile GenshinDataFileName = "ReliquaryCodexExcelConfigData"

	WeaponDataFile        GenshinDataFileName = "WeaponExcelConfigData"
	WeaponAscensionFile   GenshinDataFileName = "WeaponPromoteExcelConfigData"
	WeaponStatCurveFile   GenshinDataFileName = "WeaponCurveExcelConfigData"
	WeaponReleaseInfoFile GenshinDataFileName = "WeaponCodexExcelConfigData"

	TextMapFile           GenshinDataFileName = "ManualTextMapConfigData"
	TravelerDataFile      GenshinDataFileName = "AvatarHeroEntityExcelConfigData"
	ArchonDataFile        GenshinDataFileName = "TrialAvatarFetterDataConfigData"
	MaterialDataFile      GenshinDataFileName = "MaterialExcelConfigData"
	FriendshipRewardFile  GenshinDataFileName = "FetterCharacterCardExcelConfigData"
	RewardDataFile        GenshinDataFileName = "RewardExcelConfigData"
	ProfilePictureFile    GenshinDataFileName = "ProfilePictureExcelConfigData"
	TheaterDifficultyFile GenshinDataFileName = "RoleCombatDifficultyExcelConfigData"
)

// IsValid checks if the filename is a valid Excel config data file
func (f GenshinDataFileName) IsValid() bool {
	switch f {
	case CharacterDataFile, CharacterProfileFile, CharacterCostumeFile,
		CharacterSkillDepotFile, CharacterSkillFile, CharacterTalentFile,
		CharacterConstellationFile, CharacterAscensionFile, CharacterStatCurveFile,
		CharacterReleaseInfoFile, WeaponDataFile, WeaponAscensionFile,
		WeaponStatCurveFile, WeaponReleaseInfoFile, ArtifactSetBonusFile,
		ArtifactDataFile, ArtifactMainStatFile, ArtifactSubStatFile,
		ArtifactSetDataFile, ArtifactRarityDataFile, TextMapFile,
		TravelerDataFile, ArchonDataFile, MaterialDataFile,
		FriendshipRewardFile, RewardDataFile, ProfilePictureFile,
		TheaterDifficultyFile:
		return true
	}
	return false
}

// String returns the string representation of the filename
func (f GenshinDataFileName) String() string {
	return string(f)
}

func GetGenshinDataFileNames() []GenshinDataFileName {
	return []GenshinDataFileName{
		CharacterDataFile,
		CharacterProfileFile,
		CharacterCostumeFile,
		CharacterSkillDepotFile,
		CharacterSkillFile,
		CharacterTalentFile,
		CharacterConstellationFile,
		CharacterAscensionFile,
		CharacterStatCurveFile,
		CharacterReleaseInfoFile,
		WeaponDataFile,
		WeaponAscensionFile,
		WeaponStatCurveFile,
		WeaponReleaseInfoFile,
		ArtifactSetBonusFile,
		ArtifactDataFile,
		ArtifactMainStatFile,
		ArtifactSubStatFile,
		ArtifactSetDataFile,
		ArtifactRarityDataFile,
		TextMapFile,
		TravelerDataFile,
		ArchonDataFile,
		MaterialDataFile,
		FriendshipRewardFile,
		RewardDataFile,
		ProfilePictureFile,
		TheaterDifficultyFile,
	}
}
