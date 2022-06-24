package mainpac

import "time"

func (s *Service) GoYouTube() {
	for {
		s.YouTubeCheck()

		s.YouTube.NumberIterations++
		s.YouTube.LastTime = time.Now()
		time.Sleep(time.Minute * time.Duration(s.YouTube.CycleDurationMinutes))

		// Pause
		if s.YouTube.Pause == 1 {
			s.YouTube.Pause = 2
			<-s.YouTube.PauseWaitChannel
		}
	}
}

func (s *Service) YouTubeCheck() {
	//TODO
}

func (s *Service) YouTubePause() {
	switch s.YouTube.Pause {
	case 0:
		s.YouTube.Pause = 1
	case 1:
		s.YouTube.Pause = 0
	case 2:
		s.YouTube.Pause = 0
		s.YouTube.PauseWaitChannel <- struct{}{}
	}
}
