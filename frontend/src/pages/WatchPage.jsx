import { useEffect, useRef, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getItemById } from "../data/catalog";
import { getSelectedProfileId } from "../api/profileApi";
import { getProgress, updateProgress, markCompleted } from "../api/playbackApi";

function formatTime(seconds) {
  if (!Number.isFinite(seconds)) return "0:00";

  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);

  return `${mins}:${String(secs).padStart(2, "0")}`;
}

export default function WatchPage() {
  const { id } = useParams();
  const videoRef = useRef(null);

  const item = getItemById(id);
  const profileId = getSelectedProfileId();

  const [showSkipIntro, setShowSkipIntro] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [resumeLoaded, setResumeLoaded] = useState(false);

  useEffect(() => {
    async function loadResume() {
      if (!item || !videoRef.current) return;

      const progress = await getProgress({
        profileId,
        contentId: item.id,
      });

      if (progress.currentPositionSec > 0 && !progress.completed) {
        videoRef.current.currentTime = progress.currentPositionSec;
        setCurrentTime(progress.currentPositionSec);
      }

      setResumeLoaded(true);
    }

    loadResume();
  }, [item, profileId]);

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

  const handleTimeUpdate = async () => {
    if (!videoRef.current) return;

    const time = Math.floor(videoRef.current.currentTime);
    setCurrentTime(time);

    const insideIntro =
      time >= item.introStartSec && time <= item.introEndSec;

    setShowSkipIntro(insideIntro);

    if (time % 10 === 0) {
      await updateProgress({
        profileId,
        contentId: item.id,
        currentPositionSec: time,
      });
    }
  };

  const handlePause = async () => {
    if (!videoRef.current) return;

    await updateProgress({
      profileId,
      contentId: item.id,
      currentPositionSec: Math.floor(videoRef.current.currentTime),
    });
  };

  const handleEnded = async () => {
    await markCompleted({
      profileId,
      contentId: item.id,
    });
  };

  const skipIntro = async () => {
    if (!videoRef.current) return;

    videoRef.current.currentTime = item.introEndSec;
    setShowSkipIntro(false);

    await updateProgress({
      profileId,
      contentId: item.id,
      currentPositionSec: item.introEndSec,
    });
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
            Current time: {formatTime(currentTime)} / {formatTime(item.durationSec)}
          </p>
        </div>
      </section>

      <section className="video-shell">
        {!resumeLoaded ? (
          <div className="video-loading">Loading resume position...</div>
        ) : null}

        <video
          ref={videoRef}
          src={item.videoUrl}
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

      <section className="watch-info-panel">
        <h2>Playback data</h2>

        <div className="watch-info-grid">
          <div>
            <span>Profile ID</span>
            <strong>{profileId}</strong>
          </div>

          <div>
            <span>Content ID</span>
            <strong>{item.id}</strong>
          </div>

          <div>
            <span>Intro</span>
            <strong>
              {formatTime(item.introStartSec)} - {formatTime(item.introEndSec)}
            </strong>
          </div>

          <div>
            <span>Status</span>
            <strong>{item.status}</strong>
          </div>
        </div>
      </section>
    </main>
  );
}