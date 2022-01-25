package utils

import pb "github.com/daniarmas/api_go/pkg"

func ParsePlatformType(tp *string) *pb.PlatformType {
	if *tp == "IOS" {
		return pb.PlatformType_IOS.Enum()
	} else if *tp == "Android" {
		return pb.PlatformType_Android.Enum()
	} else {
		return pb.PlatformType_PlatformTypeUnspecified.Enum()
	}
}

func ParseAppType(tp *string) *pb.AppType {
	if *tp == "App" {
		return pb.AppType_App.Enum()
	} else if *tp == "BusinessApp" {
		return pb.AppType_BusinessApp.Enum()
	} else {
		return pb.AppType_AppTypeUnspecified.Enum()
	}
}

func ParseItemStatusType(tp *string) *pb.ItemStatusType {
	if *tp == "Available" {
		return pb.ItemStatusType_Available.Enum()
	} else if *tp == "Unavailable" {
		return pb.ItemStatusType_Unavailable.Enum()
	} else if *tp == "Deprecated" {
		return pb.ItemStatusType_Deprecated.Enum()
	} else {
		return pb.ItemStatusType_ItemStatusTypeUnspecified.Enum()
	}
}

func ParseSearchMunicipalityType(tp string) *pb.SearchMunicipalityType {
	if tp == "More" {
		return pb.SearchMunicipalityType_More.Enum()
	} else if tp == "NoMore" {
		return pb.SearchMunicipalityType_NoMore.Enum()
	} else {
		return pb.SearchMunicipalityType_SearchMunicipalityTypeUnspecified.Enum()
	}
}

func ParseResidenceType(tp string) *pb.UserAddress_UserAddressType {
	if tp == "House" {
		return pb.UserAddress_HOUSE.Enum()
	} else if tp == "Apartament" {
		return pb.UserAddress_APARTAMENT.Enum()
	} else {
		return pb.UserAddress_UNSPECIFIED.Enum()
	}
}
