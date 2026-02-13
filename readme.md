Handler/controller → Service → Repository → Database# meeting-service

fly apps create meeting-service
fly postgres create --name meeting-service-db --region bom
fly postgres attach meeting-service-db --app meeting-service
fly deploy