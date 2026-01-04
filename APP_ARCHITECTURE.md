# Music Dating App - Architecture & Vision

**Slogan:** *Mix. Match. Love.*

**Concept:** The love child of Spotify and Tinder - meet people with the same musical taste and find love through shared music.

---

## Table of Contents
1. [App Vision & Features](#app-vision--features)
2. [Current State vs Future Vision](#current-state-vs-future-vision)
3. [High-Level Architecture](#high-level-architecture)
4. [Database Schema (Full Vision)](#database-schema-full-vision)
5. [Feature Flows](#feature-flows)
6. [Current Directory Structure](#current-directory-structure)
7. [Implementation Roadmap](#implementation-roadmap)
8. [Learning Path](#learning-path)

---

## App Vision & Features

### Core Concept
Users discover others through shared music preferences in two types of rooms:
- **Public Rooms**: Join any room based on artist preference, socialize with many people
- **Private Rooms**: Invite-only rooms for deeper connections (2+ people)

### User Journey

```
1. Sign Up/Login
   ↓
2. Select Favorite Genres & Artists
   ↓
3. Choose: Create Room OR Join Existing Room
   ↓
4. ┌─────────────────┬──────────────────┐
   │  PUBLIC ROOM    │   PRIVATE ROOM   │
   ├─────────────────┼──────────────────┤
   │ • Auto-playing  │ • Full controls  │
   │   playlist      │   (play/pause)   │
   │ • Vote for next │ • Vote or waive  │
   │   song          │   to host DJ     │
   │ • Text chat     │ • Text/Voice/    │
   │ • Optional VC   │   Video call     │
   │ • Gender colors │ • Mute options   │
   │   (M=blue,      │ • Invite only    │
   │    F=pink)      │                  │
   └─────────────────┴──────────────────┘
   ↓
5. Connect → Create Private Room → Find Love ❤️
```

---

## Feature Breakdown

### 1. Artist & Genre Selection (Pre-Room)

**Requirements:**
- User must select genres they like (e.g., Pop, Rock, Hip-Hop, R&B)
- User must select favorite artists before creating/joining rooms
- This data drives room discovery (find rooms with artists you like)

**Data Needed:**
- `genres` table
- `artists` table
- `user_genres` (M:N relationship)
- `user_artists` (M:N relationship)

---

### 2. Public Rooms

**Features:**
```
┌───────────────────────────────────────────────────┐
│          PUBLIC ROOM: "Drake Fans"                │
├───────────────────────────────────────────────────┤
│  🎵 Now Playing: "God's Plan"            [3:45]   │
│  🗳️  Vote for next song:                          │
│     [ ] "Hotline Bling" (5 votes)                 │
│     [✓] "One Dance" (8 votes)                     │
│     [ ] "In My Feelings" (3 votes)                │
│                                                   │
│  👥 Users in Room (12):                           │
│     alice (F) 🎤                                  │
│     bob (M) 🔇                                    │
│     charlie (M) 🎤                                │
│                                                   │
│  💬 Text Chat:                                    │
│     alice: This song is fire! 🔥                  │
│     charlie: Anyone else going to his concert?    │
│                                                   │
│  🎙️ Voice Chat: [Host Enabled] [Mute] [Deafen]   │
└───────────────────────────────────────────────────┘
```

**Rules:**
- ✅ Music auto-plays on room creation
- ✅ No skip/pause/stop controls
- ✅ Next song determined by timed voting
- ✅ Text chat always available
- ✅ Voice chat optional (host decides)
- ✅ Gender-coded names (blue = male, pink = female)
- ✅ Users can mute individuals or deafen all

**Data Needed:**
- `rooms` table (add: `is_public`, `artist_id`, `allow_voice_chat`)
- `playlists` table (linked to artist)
- `songs` table
- `playlist_songs` (M:N)
- `room_song_votes` table (user votes for next song)
- `chat_messages` table
- `room_members` table (already exists)

---

### 3. Private Rooms

**Features:**
```
┌───────────────────────────────────────────────────┐
│        PRIVATE ROOM: "Alice & Bob"                │
├───────────────────────────────────────────────────┤
│  🎵 Now Playing: "Perfect" - Ed Sheeran           │
│  [⏮️] [⏸️] [⏭️] [🔀] [🔁]                         │
│                                                   │
│  👥 Members (2):                                  │
│     alice (F) 📹                                  │
│     bob (M) 📹                                    │
│                                                   │
│  🎙️ Video Call Active                            │
│  ┌─────────────┐  ┌─────────────┐               │
│  │   Alice     │  │    Bob      │               │
│  │   [Video]   │  │   [Video]   │               │
│  └─────────────┘  └─────────────┘               │
│                                                   │
│  💬 Chat:                                         │
│     alice: I love this song ❤️                    │
│     bob: Me too! Want to meet up?                │
└───────────────────────────────────────────────────┘
```

**Rules:**
- ✅ Minimum 2 people to create
- ✅ Unlimited max capacity
- ✅ Full music controls (play, pause, stop, skip, rewind)
- ✅ Voting still available OR users can waive votes to let host DJ
- ✅ Text, voice, or video chat
- ✅ Invite-only (members must be added)

**Data Needed:**
- `rooms` table (add: `is_public`, `host_user_id`, `voting_waived`)
- `room_invites` table
- `voice_call_sessions` table
- `video_call_sessions` table

---

## Current State vs Future Vision

### ✅ What Exists Now (Current Codebase)

```
Database:
├── users (id, username, display_name, gender, created_at)
├── rooms (id, name, created_by, is_active, created_at)
└── room_members (id, room_id, user_id, joined_at, left_at)

Code Structure:
├── routes/ - Basic HTTP routing
├── handlers/ - CRUD for users, rooms, room_members
├── store/ - Database operations
├── models/ - SQLBoiler generated models
└── migration/ - 3 basic migrations
```

**Current Capabilities:**
- ✅ Create users
- ✅ Create rooms
- ✅ Join/leave rooms
- ❌ No music integration
- ❌ No artist/genre selection
- ❌ No playlists or voting
- ❌ No chat functionality
- ❌ No voice/video calls
- ❌ No public vs private room distinction

---

### 🎯 What Needs to Be Built

```
Phase 1: Music Foundation
├── Artists, Genres, Songs, Playlists
├── User preferences (favorite artists/genres)
└── Room-playlist association

Phase 2: Room Features
├── Public vs Private room types
├── Voting system for next song
├── Music playback state management
└── Host permissions

Phase 3: Communication
├── Text chat in rooms
├── Voice chat integration
├── Video call integration
└── Mute/deafen controls

Phase 4: Discovery & Matching
├── Room discovery by artist
├── User matching by music taste
├── Recommendations
└── Gender-based UI styling

Phase 5: Advanced Features
├── Playlist generation
├── Spotify API integration
├── Real-time sync (WebSockets)
└── Invite system for private rooms
```

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND CLIENT                      │
│              (React/Vue/Mobile App)                     │
│  - Music Player UI                                      │
│  - Chat Interface                                       │
│  - Video/Voice WebRTC                                   │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP / WebSocket
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   API GATEWAY / ROUTES                  │
│              (Go net/http + WebSocket)                  │
└────────────────────┬────────────────────────────────────┘
                     │
      ┌──────────────┼──────────────┬──────────────┐
      ▼              ▼              ▼              ▼
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│  User    │  │  Room    │  │  Music   │  │  Chat    │
│ Handlers │  │ Handlers │  │ Handlers │  │ Handlers │
└────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
     │             │             │             │
     └─────────────┼─────────────┼─────────────┘
                   ▼
          ┌─────────────────┐
          │   STORE LAYER   │
          │  - UserStore    │
          │  - RoomStore    │
          │  - MusicStore   │
          │  - ChatStore    │
          │  - VoteStore    │
          └────────┬────────┘
                   │
                   ▼
          ┌─────────────────┐
          │  MODELS LAYER   │
          │  (SQLBoiler)    │
          └────────┬────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│                   MySQL DATABASE                        │
│  Tables: users, rooms, artists, songs, playlists,       │
│          votes, messages, genres, etc.                  │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│              EXTERNAL SERVICES                          │
│  - Spotify API (for music metadata)                     │
│  - WebRTC Server (for video/voice calls)                │
│  - File Storage (for profile pics, audio)               │
└─────────────────────────────────────────────────────────┘
```

---

## Database Schema (Full Vision)

### Core Entities

```
┌──────────────────────────────────────────┐
│              USERS                       │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ username (UNIQUE)                        │
│ email (UNIQUE)                           │
│ password_hash                            │
│ display_name                             │
│ gender (ENUM: male, female, other)       │
│ profile_picture_url                      │
│ created_at                               │
└──────────────┬───────────────────────────┘
               │
               │ M:N
               ▼
┌──────────────────────────────────────────┐
│           USER_GENRES                    │
├──────────────────────────────────────────┤
│ user_id (FK → users)                     │
│ genre_id (FK → genres)                   │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│              GENRES                      │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ name (e.g., "Pop", "Rock", "Hip-Hop")    │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│           USER_ARTISTS                   │
├──────────────────────────────────────────┤
│ user_id (FK → users)                     │
│ artist_id (FK → artists)                 │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│              ARTISTS                     │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ name (e.g., "Drake", "Ed Sheeran")       │
│ spotify_artist_id                        │
│ genre_id (FK → genres)                   │
│ image_url                                │
└──────────────┬───────────────────────────┘
               │
               │ 1:N
               ▼
┌──────────────────────────────────────────┐
│              ROOMS                       │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ name                                     │
│ is_public (BOOLEAN)                      │
│ created_by (FK → users)                  │
│ artist_id (FK → artists)                 │
│ current_playlist_id (FK → playlists)     │
│ current_song_id (FK → songs)             │
│ allow_voice_chat (BOOLEAN)               │
│ voting_waived (BOOLEAN)                  │
│ is_active                                │
│ created_at                               │
└──────────────┬───────────────────────────┘
               │
               │ 1:N
               ▼
┌──────────────────────────────────────────┐
│         ROOM_MEMBERS                     │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ room_id (FK → rooms)                     │
│ user_id (FK → users)                     │
│ is_host (BOOLEAN)                        │
│ is_muted (BOOLEAN)                       │
│ is_deafened (BOOLEAN)                    │
│ joined_at                                │
│ left_at (NULLABLE)                       │
└──────────────────────────────────────────┘
```

### Music Entities

```
┌──────────────────────────────────────────┐
│            PLAYLISTS                     │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ name                                     │
│ artist_id (FK → artists)                 │
│ created_by (FK → users, NULLABLE)        │
│ created_at                               │
└──────────────┬───────────────────────────┘
               │
               │ M:N
               ▼
┌──────────────────────────────────────────┐
│         PLAYLIST_SONGS                   │
├──────────────────────────────────────────┤
│ playlist_id (FK → playlists)             │
│ song_id (FK → songs)                     │
│ position (INT)                           │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│              SONGS                       │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ title                                    │
│ artist_id (FK → artists)                 │
│ spotify_track_id                         │
│ duration_ms                              │
│ album_name                               │
│ album_cover_url                          │
│ preview_url                              │
└──────────────────────────────────────────┘
```

### Interaction Entities

```
┌──────────────────────────────────────────┐
│          ROOM_SONG_VOTES                 │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ room_id (FK → rooms)                     │
│ user_id (FK → users)                     │
│ song_id (FK → songs)                     │
│ vote_session_id (UUID)                   │
│ created_at                               │
│ UNIQUE(room_id, user_id, vote_session)   │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│          CHAT_MESSAGES                   │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ room_id (FK → rooms)                     │
│ user_id (FK → users)                     │
│ message_text                             │
│ created_at                               │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│          ROOM_INVITES                    │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ room_id (FK → rooms)                     │
│ invited_by (FK → users)                  │
│ invited_user_id (FK → users)             │
│ status (ENUM: pending, accepted, rejected)│
│ created_at                               │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│       VOICE_CALL_SESSIONS                │
├──────────────────────────────────────────┤
│ id (PK)                                  │
│ room_id (FK → rooms)                     │
│ webrtc_session_id                        │
│ started_at                               │
│ ended_at (NULLABLE)                      │
└──────────────────────────────────────────┘
```

---

## Feature Flows

### Flow 1: Creating a Public Room

```
User: "I want to create a room for Drake fans"

Step 1: Select Artist
┌─────────────────────────┐
│ Select Artist:          │
│ [Search: Drake_____]    │
│                         │
│ Results:                │
│ ✓ Drake                 │
│   Drake Bell            │
│   Nick Drake            │
└─────────────────────────┘

Step 2: Create Room
┌─────────────────────────┐
│ Room Name: Drake Fans   │
│ Type: ○ Public          │
│       ○ Private         │
│ Allow Voice: ☑          │
│ [Create Room]           │
└─────────────────────────┘

Step 3: Room Created → Music Auto-Plays
┌─────────────────────────────────────┐
│ 🎵 Now Playing: "God's Plan"        │
│ 🗳️  Vote for next song...           │
│ 👥 1 user in room (you)             │
│ 💬 Chat: Room created! Say hi!      │
└─────────────────────────────────────┘
```

**Backend Flow:**
```
POST /rooms
{
  "name": "Drake Fans",
  "is_public": true,
  "artist_id": 42,
  "allow_voice_chat": true
}

1. Create room in DB
2. Fetch artist's default playlist
3. Add creator as room_member (host)
4. Start playback of first song
5. WebSocket: Broadcast room state to all clients
```

---

### Flow 2: Voting for Next Song (Public Room)

```
Timer: "30 seconds left in current song"

Server initiates vote:
┌─────────────────────────────────────┐
│ 🗳️  Vote for next song:             │
│ [ ] "Hotline Bling" (0 votes)       │
│ [ ] "One Dance" (0 votes)           │
│ [ ] "In My Feelings" (0 votes)      │
│                                     │
│ Time remaining: 25s                 │
└─────────────────────────────────────┘

Users vote:
POST /rooms/{id}/vote
{ "song_id": 123 }

After voting closes:
┌─────────────────────────────────────┐
│ 🎵 Now Playing: "One Dance"         │
│    (Won with 8 votes!)              │
└─────────────────────────────────────┘
```

**Backend Flow:**
```
1. When song reaches 80% completion:
   - Generate vote_session_id (UUID)
   - Select 3 random songs from playlist
   - Broadcast vote options via WebSocket

2. Users submit votes:
   - Insert into room_song_votes
   - Enforce UNIQUE constraint (1 vote per user per session)

3. When timer expires:
   - Count votes: SELECT song_id, COUNT(*) FROM room_song_votes WHERE vote_session_id = ? GROUP BY song_id
   - Pick winner
   - Update room.current_song_id
   - Broadcast song change via WebSocket
```

---

### Flow 3: Creating a Private Room

```
User: "I want to chat privately with alice"

Step 1: Create Private Room
┌─────────────────────────┐
│ Room Name: Me & Alice   │
│ Type: ○ Public          │
│       ● Private         │
│ Invite:                 │
│   [alice ✓]             │
│   [Search users...]     │
│ [Create Room]           │
└─────────────────────────┘

Step 2: Alice Receives Invite
┌─────────────────────────┐
│ 🔔 New Invite!          │
│ bob invited you to      │
│ "Me & Alice"            │
│ [Accept] [Decline]      │
└─────────────────────────┘

Step 3: Alice Joins → Full Controls Available
┌─────────────────────────────────────┐
│ 🎵 Paused: "Perfect" - Ed Sheeran   │
│ [⏮️] [▶️] [⏭️] [🔀] [🔁]             │
│ 👥 2 members                        │
│ 💬 Chat active                      │
│ 📹 [Start Video Call]               │
└─────────────────────────────────────┘
```

**Backend Flow:**
```
POST /rooms
{
  "name": "Me & Alice",
  "is_public": false,
  "invites": [123] // user_id of alice
}

1. Create room (is_public = false)
2. Add creator as host
3. Insert invite into room_invites
4. Send notification to alice (WebSocket/push)

POST /rooms/invites/{id}/accept
1. Update invite status = 'accepted'
2. Add alice to room_members
3. Notify room members via WebSocket
```

---

## Current Directory Structure

```
mlm/
├── main.go                    # ✅ Entry point (basic)
├── go.mod                     # ✅ Dependencies
├── docker-compose.yml         # ✅ MySQL setup
│
├── routes/
│   └── routes.go              # ✅ Basic routing (needs expansion)
│
├── handlers/
│   ├── users.go               # ✅ User CRUD
│   ├── rooms.go               # ✅ Room CRUD (needs expansion)
│   ├── room_members.go        # ✅ Join/leave
│   ├── artists.go             # ❌ TODO
│   ├── music.go               # ❌ TODO
│   ├── voting.go              # ❌ TODO
│   ├── chat.go                # ❌ TODO
│   └── websocket.go           # ❌ TODO
│
├── store/
│   ├── users_store.go         # ✅ User operations
│   ├── room_store.go          # ✅ Room operations (needs expansion)
│   ├── room_members_store.go  # ✅ Membership operations
│   ├── artist_store.go        # ❌ TODO
│   ├── music_store.go         # ❌ TODO
│   ├── vote_store.go          # ❌ TODO
│   └── chat_store.go          # ❌ TODO
│
├── models/                    # ✅ SQLBoiler (needs regeneration after new migrations)
│
└── migration/
    ├── 01_create_users.up.sql       # ✅ Exists
    ├── 02_create_rooms.up.sql       # ✅ Exists (needs update)
    ├── 03_create_room_members.up.sql # ✅ Exists (needs update)
    ├── 04_create_genres.up.sql      # ❌ TODO
    ├── 05_create_artists.up.sql     # ❌ TODO
    ├── 06_create_songs.up.sql       # ❌ TODO
    ├── 07_create_playlists.up.sql   # ❌ TODO
    ├── 08_create_user_artists.up.sql # ❌ TODO
    ├── 09_create_votes.up.sql       # ❌ TODO
    ├── 10_create_chat_messages.up.sql # ❌ TODO
    └── 11_create_invites.up.sql     # ❌ TODO
```

---

## Implementation Roadmap

### Phase 1: Music Foundation (Weeks 1-2)
**Goal:** Add music entities to the database

```
Tasks:
1. ✅ Design database schema (see above)
2. ❌ Create migrations for:
   - genres
   - artists
   - songs
   - playlists
   - playlist_songs
   - user_genres
   - user_artists
3. ❌ Regenerate SQLBoiler models
4. ❌ Create store layer for music entities
5. ❌ Create handlers for:
   - GET /genres (list all genres)
   - GET /artists?genre_id=1 (list artists by genre)
   - GET /artists/{id}/songs (get artist's songs)
   - POST /users/{id}/artists (save user's favorite artists)
6. ❌ Integrate Spotify API (optional: for fetching real music data)
```

### Phase 2: Enhanced Rooms (Weeks 3-4)
**Goal:** Public vs Private rooms with playlists

```
Tasks:
1. ❌ Update rooms table:
   - Add is_public, artist_id, current_song_id, allow_voice_chat
2. ❌ Update room_members table:
   - Add is_host, is_muted, is_deafened
3. ❌ Create room handlers:
   - POST /rooms (create with artist selection)
   - GET /rooms?artist_id=1&is_public=true (discover rooms)
   - POST /rooms/{id}/start-playback
4. ❌ Implement playlist assignment to rooms
5. ❌ Create room state management (current song, playback position)
```

### Phase 3: Voting System (Week 5)
**Goal:** Democratic song selection in public rooms

```
Tasks:
1. ❌ Create migrations:
   - room_song_votes table
2. ❌ Create vote store & handlers:
   - POST /rooms/{id}/vote { song_id }
   - GET /rooms/{id}/vote-results
3. ❌ Implement vote session logic:
   - Trigger vote when song is 80% complete
   - Close vote after 30 seconds
   - Select winner and queue next song
4. ❌ Add WebSocket events for real-time vote updates
```

### Phase 4: Chat (Week 6)
**Goal:** Text chat in all rooms

```
Tasks:
1. ❌ Create migrations:
   - chat_messages table
2. ❌ Implement WebSocket for chat:
   - WS /rooms/{id}/chat
   - Broadcast messages to all room members
3. ❌ Create chat handlers:
   - POST /rooms/{id}/messages (send message)
   - GET /rooms/{id}/messages?limit=50 (fetch history)
4. ❌ Add chat store for message persistence
```

### Phase 5: Voice & Video (Weeks 7-8)
**Goal:** Real-time communication

```
Tasks:
1. ❌ Research WebRTC implementation
2. ❌ Set up signaling server (WebSocket)
3. ❌ Create handlers:
   - POST /rooms/{id}/voice-call/start
   - POST /rooms/{id}/voice-call/join
   - POST /rooms/{id}/mute
   - POST /rooms/{id}/deafen
4. ❌ Implement peer-to-peer connections
5. ❌ Add video call support
```

### Phase 6: Private Room Features (Week 9)
**Goal:** Full music controls & invites

```
Tasks:
1. ❌ Create room_invites table
2. ❌ Implement invite system:
   - POST /rooms/{id}/invite { user_id }
   - POST /rooms/invites/{id}/accept
   - POST /rooms/invites/{id}/decline
3. ❌ Add music controls for private rooms:
   - POST /rooms/{id}/play
   - POST /rooms/{id}/pause
   - POST /rooms/{id}/skip
   - POST /rooms/{id}/seek { position_ms }
4. ❌ Implement "waive votes to host" feature
```

### Phase 7: UI/UX Polish (Week 10)
**Goal:** Gender colors, room discovery, recommendations

```
Tasks:
1. ❌ Frontend: Blue/pink name colors based on gender
2. ❌ Create room discovery page:
   - Filter by artist
   - Show active users count
   - Show current song playing
3. ❌ Build recommendation engine:
   - Suggest users with similar music taste
   - Suggest rooms based on favorite artists
4. ❌ Profile pictures and user profiles
```

---

## Learning Path

Since you're building toward this vision, here's how to learn the codebase:

### Phase 1: Understand Current Foundation (Days 1-2)
```
1. Read existing migrations (users, rooms, room_members)
2. Understand current models (Users.go, Rooms.go)
3. Trace a request: POST /users → handlers → store → DB
4. Run the app locally and test CRUD endpoints
```

### Phase 2: Design Database Schema (Day 3)
```
1. Review the full schema in this doc
2. Draw ER diagrams on paper
3. Understand relationships:
   - Users ←→ Artists (M:N via user_artists)
   - Rooms → Artists (N:1)
   - Rooms ←→ Songs (via playlists)
```

### Phase 3: Build Music Entities (Days 4-7)
```
1. Write migrations for genres, artists, songs
2. Run SQLBoiler to generate models
3. Create store layer for each entity
4. Create handlers for each entity
5. Test with Postman/curl
```

### Phase 4: Implement Room Logic (Days 8-14)
```
1. Enhance rooms table with new fields
2. Implement public vs private logic
3. Add playlist selection
4. Test room creation flow
```

### Phase 5: Add Real-Time Features (Days 15-21)
```
1. Learn WebSockets in Go (gorilla/websocket)
2. Implement chat
3. Implement voting
4. Test with multiple clients
```

### Phase 6: Integrate WebRTC (Days 22-30)
```
1. Learn WebRTC basics (signaling, ICE, STUN/TURN)
2. Implement voice chat
3. Implement video chat
4. Test with real users
```

---

## Key Technologies to Learn

1. **Go Fundamentals**
   - Goroutines & channels (for WebSocket concurrency)
   - HTTP handlers
   - JSON encoding/decoding
   - Database transactions

2. **Database**
   - MySQL relationships (1:N, M:N)
   - Indexing for performance
   - UNIQUE constraints for voting

3. **SQLBoiler**
   - Code generation from schema
   - Relationships and eager loading
   - Query builders

4. **WebSockets**
   - gorilla/websocket library
   - Broadcasting to multiple clients
   - Connection management

5. **WebRTC**
   - Signaling protocols
   - NAT traversal (STUN/TURN servers)
   - Peer connections

6. **External APIs**
   - Spotify Web API (for music metadata)
   - OAuth2 for Spotify integration

7. **Frontend (if building UI)**
   - React/Vue with WebSocket client
   - Audio player libraries (Howler.js, WaveSurfer)
   - WebRTC client libraries

---

## Next Steps

1. **Review this document** - Make sure the vision aligns with your goals
2. **Set up dev environment** - Get MySQL running, test current endpoints
3. **Start Phase 1** - Create migrations for genres, artists, songs
4. **Build incrementally** - Don't try to build everything at once
5. **Test frequently** - Use Postman, write tests as you go

---

**Questions to Consider:**

1. Do you want to integrate with Spotify API or build your own music catalog?
2. Should users authenticate with Spotify accounts or separate accounts?
3. Do you need mobile apps (iOS/Android) or web-only to start?
4. What's your target launch timeline?
5. Do you need real audio playback or just metadata/voting simulation?

---

**You're building something cool!** The combination of music + social + dating is powerful. Take it step by step, and you'll get there. Good luck! 🎵❤️
