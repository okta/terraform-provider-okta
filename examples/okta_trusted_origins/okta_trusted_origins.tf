data "okta_trusted_origins" "getAllTrustedOrigins" {
    query = null
}

data "okta_trusted_origins" "getAllTrustedOriginsFilterByQuery" {
    query = {
        limit : 100,
        filter: "%28id+eq+%22tosue7JvguwJ7U6kz0g3%22+or+id+eq+%22tos10hzarOl8zfPM80g4%22%29"
    }
}