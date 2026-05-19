import { useEffect, useRef, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getContentById } from "../api/contentApi";
import { getSelectedProfileId } from "../api/profileApi";
import {
  getProgress,
  updateProgress,
  markCompleted,
  startPlayback,
} from "../api/playbackApi";

function formatTime(seconds) {
  if (!Number.isFinite(seconds)) return "0:00";

  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);

  return `${mins}:${String(secs).padStart(2, "0")}`;
}

function inferVideoUrlFromPoster(item) {
  const posterUrl = item?.posterUrl || item?.imageUrl || "";

  const match = posterUrl.match(/\/images\/(.+)\.(jpg|jpeg|png|webp)$/i);

  if (!match) {
    return "/videos/demo-episode.mp4";
  }

  return `/videos/${match[1]}.mp4`;
}

export default function WatchPage() {
  const { id } = useParams();
  const videoRef = useRef(null);
  const lastSavedSecondRef = useRef(-1);

  const profileId = getSelectedProfileId();

  const [item, setItem] = useState(null);
  const [sessionId, setSessionId] = useState("session-created");
  const [showSkipIntro, setShowSkipIntro] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [resumeLoaded, setResumeLoaded] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadItemAndResume() {
      try {
        const content = await getContentById(id);
        setItem(content);

        try {
          const startData = await startPlayback({
            profileId,
            contentId: content.id,
            episodeId: "episode-1",
            durationSec: content.durationSec || 3600,
            device: "web",
          });

          setSessionId(startData.sessionId || "session-created");
        } catch (error) {
          console.warn("Start playback failed, continuing anyway:", error);
        }

        const progress = await getProgress({
          profileId,
          contentId: content.id,
          episodeId: "episode-1",
        });

        if (progress.currentPositionSec > 0 && !progress.completed) {
          setCurrentTime(progress.currentPositionSec);
        }
      } catch (error) {
        console.error("Failed to load watch page:", error);
      } finally {
        setResumeLoaded(true);
        setLoading(false);
      }
    }

    loadItemAndResume();
  }, [id, profileId]);

  useEffect(() => {
    if (!videoRef.current || !resumeLoaded) return;

    if (currentTime > 0) {
      videoRef.current.currentTime = currentTime;
    }
  }, [resumeLoaded, currentTime]);

  if (loading) {
    return (
      <main className="not-found-page">
        <h1>Loading video...</h1>
      </main>
    );
  }

  if (!item) {
    return (
      <main className="not-found-page">
        <h1>Video not found</h1>
        <Link to="/" className="primary-btn primary-btn--inline">
          Go Home
        </Link>
      </main>
    );
  }

  const durationSec = item.durationSec || 3600;
  const introStartSec = item.introStartSec ?? 8;
  const introEndSec = item.introEndSec ?? 72;
  const videoUrl = item.videoUrl || "/videos/demo.mp4";

  const saveProgress = async (time) => {
    try {
      await updateProgress({
        profileId,
        contentId: item.id,
        episodeId: "episode-1",
        sessionId,
        currentPositionSec: time,
        durationSec,
      });
    } catch (error) {
      console.warn("Failed to save progress:", error);
    }
  };

  const handleTimeUpdate = async () => {
    if (!videoRef.current) return;

    const time = Math.floor(videoRef.current.currentTime);
    setCurrentTime(time);

    const insideIntro = time >= introStartSec && time <= introEndSec;
    setShowSkipIntro(insideIntro);

    const shouldSave =
      time > 0 &&
      time % 10 === 0 &&
      lastSavedSecondRef.current !== time;

    if (shouldSave) {
      lastSavedSecondRef.current = time;
      await saveProgress(time);
    }
  };

  const handlePause = async () => {
    if (!videoRef.current) return;

    const time = Math.floor(videoRef.current.currentTime);
    lastSavedSecondRef.current = time;

    await saveProgress(time);
  };

  const handleEnded = async () => {
    try {
      await markCompleted({
        profileId,
        contentId: item.id,
        episodeId: "episode-1",
        sessionId,
        durationSec,
      });
    } catch (error) {
      console.warn("Failed to mark completed:", error);
    }
  };

  const skipIntro = async () => {
    if (!videoRef.current) return;

    videoRef.current.currentTime = introEndSec;
    setCurrentTime(introEndSec);
    setShowSkipIntro(false);
    lastSavedSecondRef.current = introEndSec;

    await saveProgress(introEndSec);
  };

  return (
    <main className="watch-page">
      <section className="watch-header">
        <Link to={`/title/${item.id}`} className="back-link">
          ‹ Back to details
        </Link>

        <div>
          <h1>{item.title}</h1>
          <p>
            Current time: {formatTime(currentTime)} / {formatTime(durationSec)}
          </p>
        </div>
      </section>

      <section className="video-shell">
        {!resumeLoaded ? (
          <div className="video-loading">Loading resume position...</div>
        ) : null}

        <video
          ref={videoRef}
          src={videoUrl}
          controls
          className="video-player"
          onTimeUpdate={handleTimeUpdate}
          onPause={handlePause}
          onEnded={handleEnded}
        />

        {showSkipIntro ? (
          <button className="skip-intro-btn" onClick={skipIntro}>
            Skip Intro
          </button>
        ) : null}
      </section>
    </main>
  );
}