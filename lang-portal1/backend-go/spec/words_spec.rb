require 'rspec'
require 'httparty'
require 'json'

def api_url
  'http://localhost:8080/api'
end

RSpec.describe 'Words API' do
  let(:base_url) { "#{api_url}/words" }
  let(:valid_word) do
    {
      japanese: '猫',
      romaji: 'neko',
      english: 'cat',
      parts: ['noun']
    }
  end

  describe 'GET /words' do
    before do
      @response = HTTParty.get(base_url)
    end

    it 'returns 200 status code' do
      expect(@response.code).to eq(200)
    end

    it 'returns a paginated list of words' do
      expect(@response.parsed_response).to include('items', 'current_page', 'total_pages', 'total_items', 'items_per_page')
    end

    it 'returns items as an array' do
      expect(@response.parsed_response['items']).to be_an(Array)
    end
  end

  describe 'POST /words' do
    context 'with valid parameters' do
      before do
        @response = HTTParty.post(
          base_url,
          body: valid_word.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 201 status code' do
        expect(@response.code).to eq(201)
      end

      it 'returns the created word' do
        expect(@response.parsed_response).to include(
          'id',
          'japanese' => valid_word[:japanese],
          'romaji' => valid_word[:romaji],
          'english' => valid_word[:english]
        )
      end
    end

    context 'with invalid parameters' do
      before do
        @response = HTTParty.post(
          base_url,
          body: { japanese: '猫' }.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 400 status code' do
        expect(@response.code).to eq(400)
      end

      it 'returns an error message' do
        expect(@response.parsed_response).to include('error')
      end
    end
  end

  describe 'GET /words/:id' do
    context 'when word exists' do
      before do
        # Create a word first
        post_response = HTTParty.post(
          base_url,
          body: valid_word.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
        @word_id = post_response.parsed_response['id']
        @response = HTTParty.get("#{base_url}/#{@word_id}")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns the word details' do
        expect(@response.parsed_response).to include(
          'id' => @word_id,
          'japanese' => valid_word[:japanese],
          'romaji' => valid_word[:romaji],
          'english' => valid_word[:english]
        )
      end
    end

    context 'when word does not exist' do
      before do
        @response = HTTParty.get("#{base_url}/999999")
      end

      it 'returns 404 status code' do
        expect(@response.code).to eq(404)
      end

      it 'returns an error message' do
        expect(@response.parsed_response).to include('error')
      end
    end
  end

  describe 'PUT /words/:id' do
    before do
      # Create a word first
      post_response = HTTParty.post(
        base_url,
        body: valid_word.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @word_id = post_response.parsed_response['id']
    end

    context 'with valid parameters' do
      before do
        @updated_word = valid_word.merge(english: 'kitty')
        @response = HTTParty.put(
          "#{base_url}/#{@word_id}",
          body: @updated_word.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns the updated word' do
        expect(@response.parsed_response).to include(
          'id' => @word_id,
          'english' => 'kitty'
        )
      end
    end

    context 'with invalid parameters' do
      before do
        @response = HTTParty.put(
          "#{base_url}/#{@word_id}",
          body: { japanese: '' }.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 400 status code' do
        expect(@response.code).to eq(400)
      end

      it 'returns an error message' do
        expect(@response.parsed_response).to include('error')
      end
    end
  end

  describe 'DELETE /words/:id' do
    context 'when word exists' do
      before do
        # Create a word first
        post_response = HTTParty.post(
          base_url,
          body: valid_word.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
        @word_id = post_response.parsed_response['id']
        @response = HTTParty.delete("#{base_url}/#{@word_id}")
      end

      it 'returns 204 status code' do
        expect(@response.code).to eq(204)
      end

      it 'removes the word' do
        get_response = HTTParty.get("#{base_url}/#{@word_id}")
        expect(get_response.code).to eq(404)
      end
    end

    context 'when word does not exist' do
      before do
        @response = HTTParty.delete("#{base_url}/999999")
      end

      it 'returns 204 status code' do
        expect(@response.code).to eq(204)
      end
    end
  end
end
