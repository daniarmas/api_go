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

func ParseOrderStatusType(tp *string) *pb.OrderStatusType {
	switch *tp {
	case "OrderStatusTypeStarted":
		return pb.OrderStatusType_OrderStatusTypeStarted.Enum()
	case "OrderStatusTypePending":
		return pb.OrderStatusType_OrderStatusTypePending.Enum()
	case "OrderStatusTypeRejected":
		return pb.OrderStatusType_OrderStatusTypeRejected.Enum()
	case "OrderStatusTypeApproved":
		return pb.OrderStatusType_OrderStatusTypeApproved.Enum()
	case "OrderStatusTypeReceived":
		return pb.OrderStatusType_OrderStatusTypeReceived.Enum()
	case "OrderStatusTypeCanceled":
		return pb.OrderStatusType_OrderStatusTypeCanceled.Enum()
	case "OrderStatusTypeDone":
		return pb.OrderStatusType_OrderStatusTypeDone.Enum()
	case "OrderStatusTypeExpired":
		return pb.OrderStatusType_OrderStatusTypeExpired.Enum()
	default:
		return pb.OrderStatusType_OrderStatusTypeUnspecified.Enum()
	}
}

func ParsePaymentMethodType(tp *string) *pb.PaymentMethodType {
	switch *tp {
	case "PaymentMethodTypeCash":
		return pb.PaymentMethodType_PaymentMethodTypeCash.Enum()
	case "PaymentMethodTypeEnzona":
		return pb.PaymentMethodType_PaymentMethodTypeEnzona.Enum()
	case "PaymentMethodTypeSolanaPay":
		return pb.PaymentMethodType_PaymentMethodTypeSolanaPay.Enum()
	default:
		return pb.PaymentMethodType_PaymentMethodTypeUnspecified.Enum()
	}
}

func ParsePartnerApplicationStatus(tp *string) *pb.PartnerApplicationStatus {
	switch *tp {
	case "PartnerApplicationStatusPending":
		return pb.PartnerApplicationStatus_PartnerApplicationStatusPending.Enum()
	case "PartnerApplicationStatusCanceled":
		return pb.PartnerApplicationStatus_PartnerApplicationStatusCanceled.Enum()
	case "PartnerApplicationStatusApproved":
		return pb.PartnerApplicationStatus_PartnerApplicationStatusApproved.Enum()
	case "PartnerApplicationStatusRejected":
		return pb.PartnerApplicationStatus_PartnerApplicationStatusRejected.Enum()
	default:
		return pb.PartnerApplicationStatus_PartnerApplicationStatusUnspecified.Enum().Enum()
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

func ParseOrderType(tp *string) *pb.OrderType {
	switch *tp {
	case "OrderTypePickUp":
		return pb.OrderType_OrderTypePickUp.Enum()
	case "OrderTypeHomeDelivery":
		return pb.OrderType_OrderTypeHomeDelivery.Enum()
	default:
		return pb.OrderType_OrderTypeUnspecified.Enum()
	}
}
