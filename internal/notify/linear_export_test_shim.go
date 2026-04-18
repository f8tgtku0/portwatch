package notify

func NewLinearNotifierWithURL(apiKey, teamID, apiURL string) Notifier {
	return &linearNotifier{
		apiKey: apiKey,
		teamID: teamID,
		apiURL: apiURL,
		client: defaultHTTPClient(),
	}
}
