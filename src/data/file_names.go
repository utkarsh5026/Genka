package data

// FileName represents valid Excel config data filenames
type FileName string

// Character-related files
const (
	CharacterDataFile          FileName = "AvatarExcelConfigData"
	CharacterProfileFile       FileName = "FetterInfoExcelConfigData"
	CharacterCostumeFile       FileName = "AvatarCostumeExcelConfigData"
	CharacterSkillDepotFile    FileName = "AvatarSkillDepotExcelConfigData"
	CharacterSkillFile         FileName = "AvatarSkillExcelConfigData"
	CharacterTalentFile        FileName = "ProudSkillExcelConfigData"
	CharacterConstellationFile FileName = "AvatarTalentExcelConfigData"
	CharacterAscensionFile     FileName = "AvatarPromoteExcelConfigData"
	CharacterStatCurveFile     FileName = "AvatarCurveExcelConfigData"
	CharacterReleaseInfoFile   FileName = "AvatarCodexExcelConfigData"

	ArtifactSetBonusFile   FileName = "EquipAffixExcelConfigData"
	ArtifactDataFile       FileName = "ReliquaryExcelConfigData"
	ArtifactMainStatFile   FileName = "ReliquaryLevelExcelConfigData"
	ArtifactSubStatFile    FileName = "ReliquaryAffixExcelConfigData"
	ArtifactSetDataFile    FileName = "ReliquarySetExcelConfigData"
	ArtifactRarityDataFile FileName = "ReliquaryCodexExcelConfigData"

	WeaponDataFile        FileName = "WeaponExcelConfigData"
	WeaponAscensionFile   FileName = "WeaponPromoteExcelConfigData"
	WeaponStatCurveFile   FileName = "WeaponCurveExcelConfigData"
	WeaponReleaseInfoFile FileName = "WeaponCodexExcelConfigData"

	TextMapFile           FileName = "ManualTextMapConfigData"
	TravelerDataFile      FileName = "AvatarHeroEntityExcelConfigData"
	ArchonDataFile        FileName = "TrialAvatarFetterDataConfigData"
	MaterialDataFile      FileName = "MaterialExcelConfigData"
	FriendshipRewardFile  FileName = "FetterCharacterCardExcelConfigData"
	RewardDataFile        FileName = "RewardExcelConfigData"
	ProfilePictureFile    FileName = "ProfilePictureExcelConfigData"
	TheaterDifficultyFile FileName = "RoleCombatDifficultyExcelConfigData"
)

// IsValid checks if the filename is a valid Excel config data file
func (f FileName) IsValid() bool {
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
func (f FileName) String() string {
	return string(f)
}
