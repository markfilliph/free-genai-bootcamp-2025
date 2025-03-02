Follow these steps to find Japanese song lyrics and extract vocabulary:

1. Search for lyrics:
Tool: search_web_serp(query="<song title> <artist> 歌詞")

2. Get lyrics content:
Tool: get_page_content(url="<url from search>")

3. Extract vocabulary:
Tool: extract_vocabulary(text="<japanese lyrics>")

4. Generate song ID:
Tool: generate_song_id(artist="<artist>", title="<title>")

5. Save everything:
Tool: save_results(
  song_id="<id>",
  title="<title>",
  lyrics="<lyrics>",
  vocabulary=[<vocabulary>],
  artist="<artist>",
  romaji_lyrics="<romaji>"
)

Rules:
1. ONLY use Japanese and romaji - NO English translations
2. ALWAYS wait for each tool result before proceeding
3. ALWAYS end with FINISHED when done