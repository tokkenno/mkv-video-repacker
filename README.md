# mkv-video-repacker
A simply Go program to repack MKV video files in bulk using `mkvmerge`

## Features
- Repack in bulk, using wildcards `videorepack *.mkv`
- Filter tracks by language and type (Ex: only spanish and english audio tracks)
- Select target main language: If select spanish as main language, set spanish audio and spanish forced subtitles as default. If no spanish audio, set original audio as default and enable spanish complete subtitles
- Patch flags and langs from track names: Some rippers put this information in track name in place of proper flags/langs, this tool can parse that info and set proper flags/langs. This program interprets common patterns used in track names and enables proper flags/langs accordingly.
- Specify original language: With this flags, some players will use the original audio language and complete subtitles in your preferred language if you select VOS mode.
- Rename output files using a template: The program use the metadata from the mkv file to rename the output files using a template. Ex: `{show} ({year}) - {seasonAndEpisode} - {title} [{resolution}; {video_codec}].mkv`

## Status
This program is in early development. I'm hardcoding some things for my use case. If someone has interest in this project, please open an issue or a PR to request features or report bugs.